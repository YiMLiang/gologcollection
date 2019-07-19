package main

import "fmt"

type Stu struct {
	Age int
	Name string
}

func (stu *Stu) setName(Name string) *Stu{
	stu.Name = Name
	return stu
}

func (stu *Stu) setAge(Age int) *Stu{
	stu.Age =Age
	return stu
}

func (stu *Stu)print(){
	fmt.Printf("Name = :%s, Age = :%d",stu.Name,stu.Age)
}

func main() {
	//链式结构
	stu := &Stu{}
	stu.setName("梁非凡").setAge(12).print()
}
