package pworm

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/anCreny/pw-orm/helpers"
)

type SVariable struct {
	o      *Operator
	pwName string // ${pwName} = Get-...Command
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

	v := &SVariable{
		o:      c.o,
		pwName: helpers.GenerateRandomString(10),
	}

	if result, err := c.o.RawCommand(fmt.Sprintf("%s = $global:Output_%s.Clone()", v.PW(), scopeID)).Run(); err != nil || result.Error() != nil {
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

// path - путь вложения до необходимого поля, пример: RecordData.IPv4Address
func (v *SVariable) TryGet(path string) (any, error) {
	if !validateVaribalePath(path) {
		return nil, fmt.Errorf("некорректный путь: %s", path)
	}

	if v.pwName == "" {
		return nil, fmt.Errorf("переменная не была создана")
	}

	var res any

	result, err := v.o.NewCommandBuilder(fmt.Sprintf("%s.%s", v.PW(), path)).Build().Run()
	if err != nil {
		return nil, fmt.Errorf("ошибка при старте команды на получение значения поля: %v", err)
	}

	if result.Error() != nil {
		return nil, result.Error()
	}

	if result.Output() == nil {
		return nil, fmt.Errorf("поле пустое")
	}

	if err := json.Unmarshal(result.Output(), &res); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании поля: %v", err)
	}

	return res, nil
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
		return nil, fmt.Errorf("произошла ошибка при создании клона переменной в сессии: %s", err)
	}

	return cloneVariable, nil
}

func validateVaribalePath(path string) bool {
	rx := regexp.MustCompile(`^\w+(\.\w+)*$`)
	return rx.MatchString(path)
}
