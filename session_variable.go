package pworm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/anCreny/pw-orm/helpers"
)

type SVariable struct {
	o        *Operator
	pwName   string // ${pwName} = Get-...Command
	cimClass map[string]any
}

func (c *SCommand) ToVariable() (*SVariable, error) {
	res, err := c.c.Run()
	if err != nil {
		return nil, err
	}

	if res.Error() != nil {
		return nil, res.Error()
	}

	if res.Output() == nil {
		return nil, fmt.Errorf("Команда не вернула данные для преобразования в переменную.")
	}

	var data map[string]any
	err = json.Unmarshal(res.Output(), &data)
	if err != nil {
		return nil, err
	}

	v := &SVariable{
		o:        c.o,
		pwName:   helpers.GenerateRandomString(10),
		cimClass: data,
	}

	if result, err := c.o.RawCommand(fmt.Sprintf("$global:%s['%s'] = $global:Output_%s.Clone()", c.o.s.ID, v.pwName, scopeID)).Run(); err != nil || result.Error() != nil {
		if err == nil {
			err = result.Error()
		}
		return nil, fmt.Errorf("Произошла ошибка при создании переменной в сессии: %s", err)
	}

	return v, nil
}

func (v *SVariable) PW() string {
	return fmt.Sprintf("$global:%s['%s']", v.o.s.ID, v.pwName)
}

// path - путь вложения до необходимого поля, пример: RecordData.IPv4Address
func (v *SVariable) TryGet(path string) (any, bool) {
	if v.cimClass == nil {
		return nil, false
	}

	if !validateVaribalePath(path) {
		return nil, false
	}

	parts := strings.Split(path, ".")

	var cimField any
	var ok bool

	for _, part := range parts {

		if cimField == nil {
			cimField, ok = v.cimClass[part]
			if !ok {
				return nil, false
			}

			continue
		}

		mapCimField, ok := cimField.(map[string]any)
		if !ok {
			return nil, false
		}

		cimField, ok = mapCimField[part]
		if !ok {
			return nil, false
		}
	}

	return cimField, true
}

// returns:
//
// bool:
//
//	true - поле установлено по указанному пути
//	false - поле по указанному пути не найдено и поэтому не было установлено
//
// error:
//
// nil - поле установлено успешно
// error - произошла ошибка при установке поля
func (v *SVariable) TrySet(path string, value any) (bool, error) {

	// Проверяем валидность переданных данных

	if path == "" {
		return false, fmt.Errorf("Ошибка при установке значения: путь не может быть пустой")
	}

	if value == nil {
		return false, fmt.Errorf("Ошибка при установке значения: значение не может быть nil")
	}

	if !validateVaribalePath(path) {
		return false, fmt.Errorf("Ошибка при установке значения: неверный путь")
	}

	if v.cimClass == nil {
		return false, fmt.Errorf("Ошибка при установке значения: переменная не была создана")
	}

	if v.o == nil {
		return false, fmt.Errorf("Ошибка при установке значения: оператор не был создан")
	}

	if v.pwName == "" {
		return false, fmt.Errorf("Ошибка при установке значения: имя переменной не было создано")
	}

	// Сначала обновим поле у необходимого объекта в PW
	if result, err := v.o.RawCommand(`
		$targetFieldPath = "` + path + `"
		$object = ` + v.PW() + `
		$newValueStr = "` + fmt.Sprint(value) + `"

		# 1. Разбиваем путь по точкам
		$parts = $targetFieldPath.Split('.')

		# 2. Начинаем с самого верхнего объекта
		$currentObject = $object

		# 3. Спускаемся по цепочке вложенности, ОСТАНАВЛИВАЯСЬ на предпоследнем элементе.
		# Нам нужен объект, который НАХОДИТСЯ ВНУТРИ RecordData.
		for ($i = 0; $i -lt $parts.Count - 1; $i++) {
    		$currentObject = $currentObject.($parts[$i])
		}

		# 4. Берем имя самого последнего поля
		$finalField = $parts[-1]

		# 5. Теперь запрашиваем CimType у того объекта, внутри которого это поле лежит
		$targetType = $currentObject.CimInstanceProperties[$finalField].CimType

		if ($targetType -eq "Instance" -or $targetType -eq "Reference") {

    	# Если это сложный .NET объект (как System.Net.IPAddress)
    	# Вытаскиваем его реальное .NET имя типа из метаданных CIM
    	$realTypeName = $object.$targetFieldPath.GetType().FullName
    	$typeObject = [type]$realTypeName
    
    	# Динамически вызываем метод Parse у этого типа
    	$parsedValue = $typeObject::Parse($newValueStr)

		} else {

    	# Если это простой тип (string, int, boolean) - PowerShell сам сделает безопасное приведение типов
    	$parsedValue = $newValueStr

		}

		# Применяем полученное значение обратно в клон
		$object.$targetFieldPath = $parsedValue

		`).Run(); err != nil || result.Error() != nil {
		return false, fmt.Errorf("Ошибка при установке значения в PowerShell: %s", result.Error())
	}

	// Получим обновленную структуру из PowerShell и сохраним ее в переменную
	result, err := v.o.NewCommandBuilder("$` + v.pwName + `").Build().Run()
	if err != nil {
		return false, fmt.Errorf("Ошибка при выполнении команды на получение обновленной структуры: %v", err)
	}

	if result.Error() != nil {
		return false, fmt.Errorf("Ошибка при получение обновленной структуры: %v", result.Error())
	}

	var cimClass map[string]any

	if err := json.Unmarshal(result.Output(), &cimClass); err != nil {
		return false, fmt.Errorf("Ошибка при декодировании обновленной структуры: %v", err)
	}

	v.cimClass = cimClass

	return true, nil
}

func validateVaribalePath(path string) bool {
	rx := regexp.MustCompile(`^\w+(\.\w+)*$`)
	return rx.MatchString(path)
}
