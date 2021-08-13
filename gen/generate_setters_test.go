package main

import (
	"fmt"
	"github.com/heimdalr/dag"
	"testing"
)

/**
TODO
Что тестируем?
	1. searchTypesInFile парсит файл и добавляет сеттеры в дерево. Тестовый тип есть в дереве. Есть нужные сеттеры.
	2. getSettersForType получает из дерева сеттеры для типа
	3. getSettersFromAnnotation парсит аннотацию и возвращает правильные сеттеры
	4. getSettersTemplate возвращает ожидаемый шаблон для указанных сеттеров
	5. generateMixinFile создаёт файл с ожидаемым содержимым
*/

func TestMain_Main(t *testing.T) {
	d := dag.NewDAG()
	searchTypesInFile("generate_setters_test.go", d)

	v, _ := d.GetVertex("main.TestStruct")
	fmt.Printf("%v", v.(*CldType).setters)
}
