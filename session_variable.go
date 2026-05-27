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
		return nil, fmt.Errorf("команда не вернула данные для преобразования в переменную")
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
		return nil, fmt.Errorf("произошла ошибка при создании переменной в сессии: %s", err)
	}

	return v, nil
}

func (v *SVariable) PW() string {
	return fmt.Sprintf("$global:%s['%s']", v.o.s.ID, v.pwName)
}

// Внутри переменной находится сложный CimClass из .NET,
// с которым нормально умеет работать только PowerShell.
// Для необходимости достать какое-то поле у переменной
// лучше воспользоваться командой operator.NewCommandBuilder(variable.PW()).Select("RecordData.IPv4Address.IPv4AddressToString", "IP").Build()
// и дальше скастить результат к необходимой структуре или еще раз
// вызвать TryGet
//
// path - путь вложения до необходимого поля, пример: RecordData.IPv4Address
func (v *SVariable) TryGet(path string) (any, error) {
	if v.cimClass == nil {
		return nil, fmt.Errorf("переменная не содержит данные")
	}

	if !validateVaribalePath(path) {
		return nil, fmt.Errorf("некорректный путь: %s", path)
	}

	parts := strings.Split(path, ".")

	var cimField any
	var ok bool

	for _, part := range parts {

		if cimField == nil {
			cimField, ok = v.cimClass[part]
			if !ok {
				return nil, fmt.Errorf("поле %s не найдено", part)
			}

			continue
		}

		mapCimField, ok := cimField.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("поле %s не является словарем", part)
		}

		cimField, ok = mapCimField[part]
		if !ok {
			return nil, fmt.Errorf("поле %s не найдено", part)
		}
	}

	return cimField, nil
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
		return false, fmt.Errorf("ошибка при установке значения: путь не может быть пустой")
	}

	if value == nil {
		return false, fmt.Errorf("ошибка при установке значения: значение не может быть nil")
	}

	if !validateVaribalePath(path) {
		return false, fmt.Errorf("ошибка при установке значения: неверный путь")
	}

	if v.cimClass == nil {
		return false, fmt.Errorf("ошибка при установке значения: переменная не была создана")
	}

	if v.o == nil {
		return false, fmt.Errorf("ошибка при установке значения: оператор не был создан")
	}

	if v.pwName == "" {
		return false, fmt.Errorf("ошибка при установке значения: имя переменной не было создано")
	}

	// Сначала обновим поле у необходимого объекта в PW
	if result, err := v.o.RawCommand(`
		$targetFieldPath = "` + path + `"
		$object = ` + v.PW() + `
		$newValueStr = "` + fmt.Sprint(value) + `"

		$parts = $targetFieldPath.Split('.')

		$currentObject = $object

		for ($i = 0; $i -lt $parts.Count - 1; $i++) {
    		$currentObject = $currentObject.($parts[$i])
		}

		$finalField = $parts[-1]

		$targetType = $currentObject.CimInstanceProperties[$finalField].CimType

		if ($targetType -eq "Instance" -or $targetType -eq "Reference") {

    	$realTypeName = $currentObject.$finalField.GetType().FullName
    	$typeObject = [type]$realTypeName
    
    	$parsedValue = $typeObject::Parse($newValueStr)

		} else {

    	$parsedValue = $newValueStr

		}

		$currentObject.$finalField = $parsedValue

		`).Run(); err != nil || result.Error() != nil {
		return false, fmt.Errorf("ошибка при установке значения в PowerShell: %s", result.Error())
	}

	// Получим обновленную структуру из PowerShell и сохраним ее в переменную
	result, err := v.o.NewCommandBuilder(v.PW()).Build().Run()
	if err != nil {
		return false, fmt.Errorf("ошибка при выполнении команды на получение обновленной структуры: %v", err)
	}

	if result.Error() != nil {
		return false, fmt.Errorf("ошибка при получение обновленной структуры: %v", result.Error())
	}

	var cimClass map[string]any

	if err := json.Unmarshal(result.Output(), &cimClass); err != nil {
		return false, fmt.Errorf("ошибка при декодировании обновленной структуры: %v", err)
	}

	v.cimClass = cimClass

	return true, nil
}

func (v *SVariable) Clone() (*SVariable, error) {

	cloneVariable := &SVariable{
		o:      v.o,
		pwName: helpers.GenerateRandomString(10),
	}

	if result, err := v.o.RawCommand(fmt.Sprintf("%s = %s.Clone()", cloneVariable.PW(), v.PW())).Run(); err != nil || result.Error() != nil {
		if err == nil {
			err = result.Error()
		}
		return nil, fmt.Errorf("произошла ошибка при создании переменной в сессии: %s", err)
	}

	result, err := v.o.NewCommandBuilder(cloneVariable.PW()).Build().Run()
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении команды на получение обновленной структуры: %v", err)
	}

	if result.Error() != nil {
		return nil, fmt.Errorf("ошибка при получение обновленной структуры: %v", result.Error())
	}

	var cimClass map[string]any

	if err := json.Unmarshal(result.Output(), &cimClass); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании обновленной структуры: %v", err)
	}

	cloneVariable.cimClass = cimClass

	return cloneVariable, nil
}

func validateVaribalePath(path string) bool {
	rx := regexp.MustCompile(`^\w+(\.\w+)*$`)
	return rx.MatchString(path)
}
