package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"producerPy/deepcopy"
	"strings"

	"go.uber.org/zap"
)

/*
object 第一级的属性覆盖
覆盖 string，int 等类型
覆盖 object类型
覆盖 array 类型。

a.b.c 溯源覆盖
*/

/*
obj.props[name] = 1
obj.props[objName] = convert2Obj("{}") // 必须是json
obj.props[arrName] = convert("[{},{}]")
obj.props[a].props[b].props.[c] = convert(input)

convert(type, input) {
	边界 不会有边界，yaml输入必须和api里面的层级定义相同
	case obj: to map[string]convert()
	case arr: to []interface{}, convert(array.elm.type, interface{}) append to arrType
	case str,bool,: to juti 类型
}
*/

func SetProp(root Prop, key string, value interface{}) error {
	_, err := SetPropAndGetOld(root, key, value)
	return err
}

func SetPropAndGetOld(root Prop, key string, value interface{}) (Prop, error) {
	_, target, err := findProp(root, key, value)
	if err != nil {
		return target, err
	}
	err = overwrite(target, value)
	return target, err
}

func overwrite(propInter Prop, value interface{}) error {
	getObjValue := func(s string) (map[string]interface{}, error) {
		immediateChildren := map[string]interface{}{}
		err := json.Unmarshal([]byte(s), &immediateChildren)
		return immediateChildren, err
	}
	switch currentProp := propInter.(type) {
	case *Object:
		v, ok := value.(string)
		if !ok {
			return errors.New(fmt.Sprintf("object string required, but:%s", value))
		}
		immediateChildren, err := getObjValue(v)
		if err != nil {
			zap.L().Error(`get object value failed`, zap.Error(err))
			return err
		}
		for childName, child := range immediateChildren {
			child2Overwrite, exists := currentProp.Props[childName]
			if !exists {
				zap.L().Error(`property to overwrite not found`, zap.String(`name`, childName), zap.Any(`input`, immediateChildren))
				return errors.New(childName + "not found")
			}
			if err := overwrite(child2Overwrite, child); err != nil {
				return err
			}
		}
		return nil
	case *Array:
		switch v := value.(type) {
		case string:
			currentProp.PropsArr[0] = NewString(v)
		case []interface{}:
			children := v
			newArray := make([]Prop, 0, len(children)) // 测试用例里面数组元素个数可能超过1个，而api解析出来的Array类型默认只有一个Prop，所以需要重新构造Prop，并且重写它们
			for _, child := range children {
				newProp := deepcopy.Copy(currentProp.PropsArr[0]).(Prop)
				if err := overwrite(newProp, child); err != nil {
					return err
				}
				newArray = append(newArray, newProp)
			}
			currentProp.PropsArr = newArray
		default:
			otherProp := NewBasicProp()
			otherProp.SetValue(v)
			currentProp.PropsArr[0] = otherProp
		}
		return nil
	case *String:
		currentProp.SetValue(value)
		return nil
	case *Bool:
		currentProp.SetValue(value)
		return nil
	case *Int:
		currentProp.SetValue(value)
		return nil
	case *Number:
		currentProp.SetValue(value)
		return nil
	default:
		panic(`type not support yet`)
	}
}

func findProp(prop Prop, key string, value interface{}) (father, target Prop, err error) {
	cur, restKey := resolveKey(key)
	target, err = prop.GetProp(cur)
	if err != nil {
		return father, nil, err
	}
	isThisLayerProp := restKey == ""
	if isThisLayerProp {
		father = prop
		return
	}

	return findProp(target, restKey, value)
}

func resolveKey(key string) (current, next string) { // key: a.b.c current: a, next : b.c
	names := strings.Split(key, ".")
	current = names[0]
	if len(names) > 0 {
		next = strings.Join(names[1:], ".")
	}
	return
}
