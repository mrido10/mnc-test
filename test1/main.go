package main

import (
	"fmt"
	"time"
)

type Q2 struct {
	Price int
	Pay   int
}

type Q4 struct {
	TotalLeave    int
	LeaveDuration int
	JoinDateStr   string
	LeavePlanStr  string
}

func main() {
	// quest 1
	quest1Input := [][]string{
		{"abcd", "acbd", "aaab", "acbd"},
		{"Satu", "Sate", "Tujuh", "Tusuk", "Tujuh", "Sate", "Bonus", "Tiga", "Puluh", "Tujuh", "Tusuk"},
		{"pisang", "goreng", "enak", "sekali", "rasanya"},
	}
	fmt.Println("Quest 1 . . . . . . ")
	for i, input := range quest1Input {
		fmt.Print(i+1, ".\t")
		num := quest1(input)
		if len(num) == 0 {
			fmt.Println("false")
			continue
		}
		fmt.Println(num)
	}

	// quest 2
	fmt.Println("\nQuest 2 . . . . . . ")
	quest2Input := []Q2{
		{Price: 700649, Pay: 800000},
		{Price: 575650, Pay: 580000},
		{Price: 657650, Pay: 600000},
	}
	for i, input := range quest2Input {
		fmt.Print(i+1, ".")
		diff := quest2(input.Price, input.Pay)
		for _, d := range diff {
			fmt.Println("\t", d)
		}
	}

	// quest 2
	fmt.Println("\nQuest 4 . . . . . . ")
	quest4Input := []Q4{
		{TotalLeave: 7, LeaveDuration: 1, JoinDateStr: "2021-05-01", LeavePlanStr: "2021-07-05"},
		{TotalLeave: 7, LeaveDuration: 3, JoinDateStr: "2021-05-01", LeavePlanStr: "2021-11-05"},
		{TotalLeave: 7, LeaveDuration: 1, JoinDateStr: "2021-01-05", LeavePlanStr: "2021-12-18"},
		{TotalLeave: 7, LeaveDuration: 3, JoinDateStr: "2021-01-05", LeavePlanStr: "2021-12-18"},
	}
	for i, input := range quest4Input {
		fmt.Print(i+1, ".\t")
		isAllow, reason := quest4(input.TotalLeave, input.LeaveDuration, input.JoinDateStr, input.LeavePlanStr)
		fmt.Println(isAllow, reason)
	}

}

func quest1(input []string) []int {
	var str string
	temp := make(map[string][]int)
	for i, v := range input {
		if str == "" && len(temp[v]) > 0 {
			str = v
		}
		temp[v] = append(temp[v], i+1)
	}
	return temp[str]
}

func quest2(price, pay int) []string {
	if pay < price {
		return []string{"kurang bayar"}
	}

	moneyExist := []int{
		100000,
		50000,
		20000,
		10000,
		5000,
		2000,
		1000,
		500,
		200,
		100,
	}

	diff := pay - price
	mod := diff % 100
	if mod > 0 {
		diff -= mod
	}

	var str []string
	for _, v := range moneyExist {
		if v > diff {
			continue
		}

		total := diff / v
		diff -= v * total

		typ := "lembar"
		if v < 1000 {
			typ = "koin"
		}
		str = append(str, fmt.Sprintf("%d %s %d", total, typ, v))

		if diff == 0 {
			break
		}
	}

	return str
}

func quest4(totalLeave, leaveDuration int, joinDateStr, leavePlanStr string) (bool, string) {
	joinDate, _ := time.Parse("2006-01-02", joinDateStr)
	leavePlan, _ := time.Parse("2006-01-02", leavePlanStr)

	leaveDateAllow := joinDate.AddDate(0, 0, 180)
	endYear := time.Date(leavePlan.Year(), time.December, 31, 0, 0, 0, 0, time.UTC)

	day := endYear.Sub(leaveDateAllow).Hours()
	leaveQuota := int((day / 24) / float64(365) * float64(totalLeave))

	if leaveQuota <= 0 {
		return false, "Belum ada quota cuti"
	}

	if leaveDuration > 3 {
		return false, "Tidak boleh melebihi 3 hari"
	}

	if leaveDuration > leaveQuota {
		return false, fmt.Sprintf("Hanya boleh mengambil %d hari cuti", leaveQuota)
	}

	if leavePlan.Before(leaveDateAllow) {
		return false, "Belum 180 hari sejak tanggal join karyawan"
	}
	return true, ""
}
