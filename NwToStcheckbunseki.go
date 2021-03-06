package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	//　ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)

	//　入力ファイル準備
	infile, err := os.Open(flag.Arg(0))
	failOnError(err)
	defer infile.Close()

	//　書き込みファイル準備
	outfile, err := os.Create("./ストレスチェック集団分析用データ.txt")
	failOnError(err)
	defer outfile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'
	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	writer.Comma = '\t'
	writer.UseCRLF = true

	log.Print("Start\r\n")
	//　タイトル行を取得
	recordHead, err := reader.Read() // 1行読み出す
	if err != io.EOF {
		failOnError(err)
	}

	// タイトル行の書き込み
	addRecordHead(&recordHead)
	writer.Write(recordHead)

	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		var out_record []string
		var A []int
		var B []int
		var C []int
		var D []int
		out_record = append(out_record, record...)
		//errPersonalInfo := record[0] + "," + record[3]

		//項目文字列を数字に変換
		for i := 23; i <= 79; i++ {
			if i <= 39 {
				setA := setAvalue(record[i])
				intA, _ := strconv.Atoi(setA)
				out_record = append(out_record, setA)
				A = append(A, intA)
			} else if i <= 68 {
				setB := setBvalue(record[i])
				intB, _ := strconv.Atoi(setB)
				out_record = append(out_record, setB)
				B = append(B, intB)
			} else if i <= 77 {
				setC := setCvalue(record[i])
				intC, _ := strconv.Atoi(setC)
				out_record = append(out_record, setC)
				C = append(C, intC)
			} else if i <= 79 {
				setD := setDvalue(record[i])
				intD, _ := strconv.Atoi(setD)
				out_record = append(out_record, setD)
				D = append(D, intD)
			}
		}
		//log.Print(errPersonalInfo + "\r\n")

		//素点換算
		var Soten []int
		Soten = append(Soten, setSoten1(record[4], A[0], A[1], A[2]))
		Soten = append(Soten, setSoten2(record[4], A[3], A[4], A[5]))
		Soten = append(Soten, setSoten3(record[4], A[6]))
		Soten = append(Soten, setSoten4(record[4], A[11], A[12], A[13]))
		Soten = append(Soten, setSoten5(record[4], A[14]))
		Soten = append(Soten, setSoten6(record[4], A[7], A[8], A[9]))
		Soten = append(Soten, setSoten7(record[4], A[10]))
		Soten = append(Soten, setSoten8(record[4], A[15]))
		Soten = append(Soten, setSoten9(record[4], A[16]))
		Soten = append(Soten, setSoten10(record[4], B[0], B[1], B[2]))
		Soten = append(Soten, setSoten11(record[4], B[3], B[4], B[5]))
		Soten = append(Soten, setSoten12(record[4], B[6], B[7], B[8]))
		Soten = append(Soten, setSoten13(record[4], B[9], B[10], B[11]))
		Soten = append(Soten, setSoten14(record[4], B[12], B[13], B[14], B[15], B[16], B[17]))
		Soten = append(Soten, setSoten15(record[4], B[18], B[19], B[20], B[21], B[22], B[23], B[24], B[25], B[26], B[27], B[28]))
		Soten = append(Soten, setSoten16(record[4], C[0], C[3], C[6]))
		Soten = append(Soten, setSoten17(record[4], C[1], C[4], C[7]))
		Soten = append(Soten, setSoten18(record[4], C[2], C[5], C[8]))
		Soten = append(Soten, setSoten19(record[4], D[0], D[1]))
		out_record = append(out_record, intToString(Soten)...)

		//素点換算_高ストレスへ転換
		var SSoten []int
		SSoten = append(SSoten, 6-Soten[0])
		SSoten = append(SSoten, 6-Soten[1])
		SSoten = append(SSoten, 6-Soten[2])
		SSoten = append(SSoten, 6-Soten[3])
		SSoten = append(SSoten, 6-Soten[4])
		SSoten = append(SSoten, Soten[5])
		SSoten = append(SSoten, Soten[6])
		SSoten = append(SSoten, Soten[7])
		SSoten = append(SSoten, Soten[8])
		SSoten = append(SSoten, Soten[9])
		SSoten = append(SSoten, 6-Soten[10])
		SSoten = append(SSoten, 6-Soten[11])
		SSoten = append(SSoten, 6-Soten[12])
		SSoten = append(SSoten, 6-Soten[13])
		SSoten = append(SSoten, 6-Soten[14])
		SSoten = append(SSoten, Soten[15])
		SSoten = append(SSoten, Soten[16])
		SSoten = append(SSoten, Soten[17])
		SSoten = append(SSoten, Soten[18])
		out_record = append(out_record, intToString(SSoten)...)

		//高ストレス計算（素点換算）
		var sotenA, sotenB, sotenC, sotenHi int
		//Ａ.心身のストレス反応 9-14
		for n := 9; n <= 14; n++ {
			sotenA = sotenA + SSoten[n]
		}
		//Ｂ.仕事のストレス要因 0-8
		for n := 0; n <= 8; n++ {
			sotenB = sotenB + SSoten[n]
		}
		//Ｃ.周囲のサポート 15-18
		for n := 15; n <= 17; n++ {
			sotenC = sotenC + SSoten[n]
		}
		//高ストレス判定（素点換算）
		if sotenA <= 12 {
			sotenHi = 1
		} else if (sotenB+sotenC) <= 26 && sotenA <= 17 {
			sotenHi = 1
		} else {
			sotenHi = 0
		}

		hiSoten := make([]int, 4)
		hiSoten[0] = sotenA
		hiSoten[1] = sotenB
		hiSoten[2] = sotenC
		hiSoten[3] = sotenHi
		out_record = append(out_record, intToString(hiSoten)...)

		//合計点数_高ストレスへ転換
		gA := A
		gB := B
		gC := C

		gA[0] = 5 - gA[0]
		gA[1] = 5 - gA[1]
		gA[2] = 5 - gA[2]
		gA[3] = 5 - gA[3]
		gA[4] = 5 - gA[4]
		gA[5] = 5 - gA[5]
		gA[6] = 5 - gA[6]
		gA[10] = 5 - gA[10]
		gA[11] = 5 - gA[11]
		gA[12] = 5 - gA[12]
		gA[14] = 5 - gA[14]
		gB[0] = 5 - gB[0]
		gB[1] = 5 - gB[1]
		gB[2] = 5 - gB[2]

		//高ストレス計算（合計点数）
		var gokeiA, gokeiB, gokeiC, gokeiHi int
		for n := range gA {
			gokeiA = gokeiA + gA[n]
		}
		for n := range gB {
			gokeiB = gokeiB + gB[n]
		}
		for n := range gC {
			gokeiC = gokeiC + gC[n]
		}

		//高ストレス判定（合計点数）
		if gokeiB >= 77 {
			gokeiHi = 1
		} else if (gokeiA+gokeiC) >= 76 && gokeiB >= 63 {
			gokeiHi = 1
		} else {
			gokeiHi = 0
		}

		hiGokei := make([]string, 4)
		hiGokei[0] = strconv.Itoa(gokeiB)
		hiGokei[1] = strconv.Itoa(gokeiA)
		hiGokei[2] = strconv.Itoa(gokeiC)
		hiGokei[3] = strconv.Itoa(gokeiHi)
		out_record = append(out_record, hiGokei...)

		//集団分析項目の文字列を数字に変換
		Ab := make([]int, 6)
		Cb := make([]int, 6)

		Ab[0] = setAbvalue(record[23])
		Ab[1] = setAbvalue(record[24])
		Ab[2] = setAbvalue(record[25])
		Ab[3] = setAbvalue(record[30])
		Ab[4] = setAbvalue(record[31])
		Ab[5] = setAbvalue(record[32])
		Cb[0] = setCbvalue(record[69])
		Cb[1] = setCbvalue(record[70])
		Cb[2] = setCbvalue(record[72])
		Cb[3] = setCbvalue(record[73])
		Cb[4] = setCbvalue(record[75])
		Cb[5] = setCbvalue(record[76])
		out_record = append(out_record, intToString(Ab)...)
		out_record = append(out_record, intToString(Cb)...)

		//集団分析　得点計算
		Sb := make([]int, 4)
		Sb[0] = Ab[0] + Ab[1] + Ab[2]
		Sb[1] = Ab[3] + Ab[4] + Ab[5]
		Sb[2] = Cb[0] + Cb[2] + Cb[4]
		Sb[3] = Cb[1] + Cb[3] + Cb[5]
		out_record = append(out_record, intToString(Sb)...)

		//集団分析　Bad
		SBad := make([]int, 19)
		SBad[0] = setBad(SSoten[0])
		SBad[1] = setBad(SSoten[1])
		SBad[2] = setBad(SSoten[2])
		SBad[3] = setBad(SSoten[3])
		SBad[4] = setBad(SSoten[4])
		SBad[5] = setBad(SSoten[5])
		SBad[6] = setBad(SSoten[6])
		SBad[7] = setBad(SSoten[7])
		SBad[8] = setBad(SSoten[8])
		SBad[9] = setBad(SSoten[15])
		SBad[10] = setBad(SSoten[16])
		SBad[11] = setBad(SSoten[17])
		SBad[12] = setBad(SSoten[18])
		SBad[13] = setBad(SSoten[9])
		SBad[14] = setBad(SSoten[10])
		SBad[15] = setBad(SSoten[11])
		SBad[16] = setBad(SSoten[12])
		SBad[17] = setBad(SSoten[13])
		SBad[18] = setBad(SSoten[14])
		out_record = append(out_record, intToString(SBad)...)

		//集団分析　Good
		SGood := make([]int, 13)
		SGood[0] = setGood(SSoten[0])
		SGood[1] = setGood(SSoten[1])
		SGood[2] = setGood(SSoten[2])
		SGood[3] = setGood(SSoten[3])
		SGood[4] = setGood(SSoten[4])
		SGood[5] = setGood(SSoten[5])
		SGood[6] = setGood(SSoten[6])
		SGood[7] = setGood(SSoten[7])
		SGood[8] = setGood(SSoten[8])
		SGood[9] = setGood(SSoten[15])
		SGood[10] = setGood(SSoten[16])
		SGood[11] = setGood(SSoten[17])
		SGood[12] = setGood(SSoten[18])
		out_record = append(out_record, intToString(SGood)...)

		//集団分析対象外
		bunsekiNo := "0"
		for _, Sotenv := range Soten {
			if Sotenv == 0 {
				bunsekiNo = "1"
			}
		}
		out_record = append(out_record, bunsekiNo)

		writer.Write(out_record)

	}

	writer.Flush()
	log.Print("Finesh !\r\n")

}

func addRecordHead(S *[]string) {
	// ファイルの先頭行に項目を追加する
	*S = append(*S, "A1")
	*S = append(*S, "A2")
	*S = append(*S, "A3")
	*S = append(*S, "A4")
	*S = append(*S, "A5")
	*S = append(*S, "A6")
	*S = append(*S, "A7")
	*S = append(*S, "A8")
	*S = append(*S, "A9")
	*S = append(*S, "A10")
	*S = append(*S, "A11")
	*S = append(*S, "A12")
	*S = append(*S, "A13")
	*S = append(*S, "A14")
	*S = append(*S, "A15")
	*S = append(*S, "A16")
	*S = append(*S, "A17")
	*S = append(*S, "B1")
	*S = append(*S, "B2")
	*S = append(*S, "B3")
	*S = append(*S, "B4")
	*S = append(*S, "B5")
	*S = append(*S, "B6")
	*S = append(*S, "B7")
	*S = append(*S, "B8")
	*S = append(*S, "B9")
	*S = append(*S, "B10")
	*S = append(*S, "B11")
	*S = append(*S, "B12")
	*S = append(*S, "B13")
	*S = append(*S, "B14")
	*S = append(*S, "B15")
	*S = append(*S, "B16")
	*S = append(*S, "B17")
	*S = append(*S, "B18")
	*S = append(*S, "B19")
	*S = append(*S, "B20")
	*S = append(*S, "B21")
	*S = append(*S, "B22")
	*S = append(*S, "B23")
	*S = append(*S, "B24")
	*S = append(*S, "B25")
	*S = append(*S, "B26")
	*S = append(*S, "B27")
	*S = append(*S, "B28")
	*S = append(*S, "B29")
	*S = append(*S, "C1")
	*S = append(*S, "C2")
	*S = append(*S, "C3")
	*S = append(*S, "C4")
	*S = append(*S, "C5")
	*S = append(*S, "C6")
	*S = append(*S, "C7")
	*S = append(*S, "C8")
	*S = append(*S, "C9")
	*S = append(*S, "D1")
	*S = append(*S, "D2")
	*S = append(*S, "量的負担")
	*S = append(*S, "質的負担")
	*S = append(*S, "身体負担")
	*S = append(*S, "対人関係")
	*S = append(*S, "職場環境")
	*S = append(*S, "コントロール")
	*S = append(*S, "技能活用")
	*S = append(*S, "適性度")
	*S = append(*S, "働き甲斐")
	*S = append(*S, "活気")
	*S = append(*S, "いらいら感")
	*S = append(*S, "疲労感")
	*S = append(*S, "不安感")
	*S = append(*S, "抑うつ感")
	*S = append(*S, "身体愁訴")
	*S = append(*S, "上司支援")
	*S = append(*S, "同僚支援")
	*S = append(*S, "家族・友人支援")
	*S = append(*S, "満足度")
	*S = append(*S, "s量的負担")
	*S = append(*S, "s質的負担")
	*S = append(*S, "s身体負担")
	*S = append(*S, "s対人関係")
	*S = append(*S, "s職場環境")
	*S = append(*S, "sコントロール")
	*S = append(*S, "s技能活用")
	*S = append(*S, "s適性度")
	*S = append(*S, "s働き甲斐")
	*S = append(*S, "s活気")
	*S = append(*S, "sいらいら感")
	*S = append(*S, "s疲労感")
	*S = append(*S, "s不安感")
	*S = append(*S, "s抑うつ感")
	*S = append(*S, "s身体愁訴")
	*S = append(*S, "s上司支援")
	*S = append(*S, "s同僚支援")
	*S = append(*S, "s家族・友人支援")
	*S = append(*S, "s満足度")
	*S = append(*S, "心身のストレス反応(素点)")
	*S = append(*S, "仕事のストレス要因(素点)")
	*S = append(*S, "周囲のサポート(素点)")
	*S = append(*S, "高ストレス(素点)")
	*S = append(*S, "心身のストレス反応(合計)")
	*S = append(*S, "仕事のストレス要因(合計)")
	*S = append(*S, "周囲のサポート(合計)")
	*S = append(*S, "高ストレス(合計)")
	*S = append(*S, "A1b")
	*S = append(*S, "A2b")
	*S = append(*S, "A3b")
	*S = append(*S, "A8b")
	*S = append(*S, "A9b")
	*S = append(*S, "A10b")
	*S = append(*S, "C1b")
	*S = append(*S, "C2b")
	*S = append(*S, "C4b")
	*S = append(*S, "C5b")
	*S = append(*S, "C7b")
	*S = append(*S, "C8b")
	*S = append(*S, "仕事の量的負担")
	*S = append(*S, "仕事のコントロール")
	*S = append(*S, "上司の支援")
	*S = append(*S, "同僚の支援")
	*S = append(*S, "量的負担_Bad")
	*S = append(*S, "質的負担_Bad")
	*S = append(*S, "身体負担_Bad")
	*S = append(*S, "対人関係_Bad")
	*S = append(*S, "職場環境_Bad")
	*S = append(*S, "コントロール_Bad")
	*S = append(*S, "技能活用_Bad")
	*S = append(*S, "適性度_Bad")
	*S = append(*S, "働き甲斐_Bad")
	*S = append(*S, "上司支援_Bad")
	*S = append(*S, "同僚支援_Bad")
	*S = append(*S, "家族・友人支援_Bad")
	*S = append(*S, "満足度_Bad")
	*S = append(*S, "活気_Bad")
	*S = append(*S, "いらいら感_Bad")
	*S = append(*S, "疲労感_Bad")
	*S = append(*S, "不安感_Bad")
	*S = append(*S, "抑うつ感_Bad")
	*S = append(*S, "身体愁訴_Bad")
	*S = append(*S, "量的負担_Good")
	*S = append(*S, "質的負担_Good")
	*S = append(*S, "身体負担_Good")
	*S = append(*S, "対人関係_Good")
	*S = append(*S, "職場環境_Good")
	*S = append(*S, "コントロール_Good")
	*S = append(*S, "技能活用_Good")
	*S = append(*S, "適性度_Good")
	*S = append(*S, "働き甲斐_Good")
	*S = append(*S, "上司支援_Good")
	*S = append(*S, "同僚支援_Good")
	*S = append(*S, "家族・友人支援_Good")
	*S = append(*S, "満足度_Good")
	*S = append(*S, "集団分析対象外")

}

func setAvalue(S string) string {

	switch S {
	case "そうだ":
		S = "1"
	case "まあそうだ":
		S = "2"
	case "ややちがう":
		S = "3"
	case "ちがう":
		S = "4"
	default:
		S = ""
	}
	return S
}

func setBvalue(S string) string {

	switch S {
	case "ほとんどなかった":
		S = "1"
	case "ときどきあった":
		S = "2"
	case "しばしばあった":
		S = "3"
	case "ほとんどいつもあった":
		S = "4"
	default:
		S = ""
	}
	return S
}

func setCvalue(S string) string {

	switch S {
	case "非常に":
		S = "1"
	case "かなり":
		S = "2"
	case "多少":
		S = "3"
	case "全くない":
		S = "4"
	default:
		S = ""
	}
	return S
}

func setDvalue(S string) string {

	switch S {
	case "満足":
		S = "1"
	case "まあ満足":
		S = "2"
	case "やや不満足":
		S = "3"
	case "不満足":
		S = "4"
	default:
		S = ""
	}
	return S
}

func setSoten1(sei string, a1, a2, a3 int) int {

	B := setBlank3(a1, a2, a3)

	// S := 15 - (a1 + a2 + a3)
	S := 15 - (B[0] + B[1] + B[2])
	if sei == "男" {
		switch S {
		case 3, 4, 5:
			S = 1
		case 6, 7:
			S = 2
		case 8, 9:
			S = 3
		case 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3, 4:
			S = 1
		case 5, 6:
			S = 2
		case 7, 8, 9:
			S = 3
		case 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten2(sei string, a4, a5, a6 int) int {

	B := setBlank3(a4, a5, a6)

	//S := 15 - (a4 + a5 + a6)
	S := 15 - (B[0] + B[1] + B[2])
	if sei == "男" {
		switch S {
		case 3, 4, 5:
			S = 1
		case 6, 7:
			S = 2
		case 8, 9:
			S = 3
		case 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3, 4:
			S = 1
		case 5, 6:
			S = 2
		case 7, 8:
			S = 3
		case 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten3(sei string, a7 int) int {

	S := 5 - a7
	if sei == "男" {
		switch S {
		case 1:
			S = 2
		case 2:
			S = 3
		case 3:
			S = 4
		case 4:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 1:
			S = 2
		case 2:
			S = 3
		case 3:
			S = 4
		case 4:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten4(sei string, a12, a13, a14 int) int {

	B := setBlank33(a12, a13, a14)

	//S := 10 - (a12 + a13) + a14
	S := 10 - (B[0] + B[1]) + B[2]
	if sei == "男" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7:
			S = 3
		case 8, 9:
			S = 4
		case 10, 11, 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7:
			S = 3
		case 8, 9:
			S = 4
		case 10, 11, 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}

	return S
}

func setSoten5(sei string, a15 int) int {

	S := 5 - a15
	if sei == "男" {
		switch S {
		case 1:
			S = 2
		case 2:
			S = 3
		case 3:
			S = 4
		case 4:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 1:
			S = 1
		case 2:
			S = 3
		case 3:
			S = 4
		case 4:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten6(sei string, a8, a9, a10 int) int {

	B := setBlank3(a8, a9, a10)

	//S := 15 - (a8 + a9 + a10)
	S := 15 - (B[0] + B[1] + B[2])
	if sei == "男" {
		switch S {
		case 3, 4:
			S = 1
		case 5, 6:
			S = 2
		case 7, 8:
			S = 3
		case 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7, 8:
			S = 3
		case 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten7(sei string, a11 int) int {

	S := a11
	if sei == "男" {
		switch S {
		case 1:
			S = 1
		case 2:
			S = 2
		case 3:
			S = 3
		case 4:
			S = 4
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 1:
			S = 1
		case 2:
			S = 2
		case 3:
			S = 3
		case 4:
			S = 4
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten8(sei string, a16 int) int {

	S := 5 - a16
	if sei == "男" {
		switch S {
		case 1:
			S = 1
		case 2:
			S = 2
		case 3:
			S = 3
		case 4:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 1:
			S = 1
		case 2:
			S = 2
		case 3:
			S = 3
		case 4:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten9(sei string, a17 int) int {

	S := 5 - a17
	if sei == "男" {
		switch S {
		case 1:
			S = 1
		case 2:
			S = 2
		case 3:
			S = 3
		case 4:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 1:
			S = 1
		case 2:
			S = 2
		case 3:
			S = 3
		case 4:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten10(sei string, b1, b2, b3 int) int {

	B := setBlank3(b1, b2, b3)

	//S := b1 + b2 + b3
	S := B[0] + B[1] + B[2]
	if sei == "男" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7:
			S = 3
		case 8, 9:
			S = 4
		case 10, 11, 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7:
			S = 3
		case 8, 9:
			S = 4
		case 10, 11, 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten11(sei string, b4, b5, b6 int) int {

	B := setBlank3(b4, b5, b6)

	//S := b4 + b5 + b6
	S := B[0] + B[1] + B[2]
	if sei == "男" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7:
			S = 3
		case 8, 9:
			S = 4
		case 10, 11, 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7, 8:
			S = 3
		case 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten12(sei string, b7, b8, b9 int) int {

	B := setBlank3(b7, b8, b9)

	//S := b7 + b8 + b9
	S := B[0] + B[1] + B[2]
	if sei == "男" {
		switch S {
		case 3:
			S = 1
		case 4:
			S = 2
		case 5, 6, 7:
			S = 3
		case 8, 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7, 8:
			S = 3
		case 9, 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten13(sei string, b10, b11, b12 int) int {

	B := setBlank3(b10, b11, b12)

	//S := b10 + b11 + b12
	S := B[0] + B[1] + B[2]
	if sei == "男" {
		switch S {
		case 3:
			S = 1
		case 4:
			S = 2
		case 5, 6, 7:
			S = 3
		case 8, 9:
			S = 4
		case 10, 11, 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3:
			S = 1
		case 4:
			S = 2
		case 5, 6, 7:
			S = 3
		case 8, 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten14(sei string, b13, b14, b15, b16, b17, b18 int) int {

	B := setBlank6(b13, b14, b15, b16, b17, b18)

	//S := b13 + b14 + b15 + b16 + b17 + b18
	S := B[0] + B[1] + B[2] + B[3] + B[4] + B[5]
	if sei == "男" {
		switch S {
		case 6:
			S = 1
		case 7, 8:
			S = 2
		case 9, 10, 11, 12:
			S = 3
		case 13, 14, 15, 16:
			S = 4
		case 17, 18, 19, 20, 21, 22, 23, 24:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 6:
			S = 1
		case 7, 8:
			S = 2
		case 9, 10, 11, 12:
			S = 3
		case 13, 14, 15, 16, 17:
			S = 4
		case 18, 19, 20, 21, 22, 23, 24:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten15(sei string, b19, b20, b21, b22, b23, b24, b25, b26, b27, b28, b29 int) int {

	B := setBlank11(b19, b20, b21, b22, b23, b24, b25, b26, b27, b28, b29)

	//S := b19 + b20 + b21 + b22 + b23 + b24 + b25 + b26 + b27 + b28 + b29
	S := B[0] + B[1] + B[2] + B[3] + B[4] + B[5] + B[6] + B[7] + B[8] + B[9] + B[10]
	if sei == "男" {
		switch {
		case S == 11:
			S = 1
		case S >= 12 && S <= 15:
			S = 2
		case S >= 16 && S <= 21:
			S = 3
		case S >= 22 && S <= 26:
			S = 4
		case S >= 27 && S <= 44:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch {
		case S >= 11 && S <= 13:
			S = 1
		case S >= 14 && S <= 17:
			S = 2
		case S >= 18 && S <= 23:
			S = 3
		case S >= 24 && S <= 29:
			S = 4
		case S >= 30 && S <= 44:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten16(sei string, c1, c4, c7 int) int {

	B := setBlank3(c1, c4, c7)

	//S := 15 - (c1 + c4 + c7)
	S := 15 - (B[0] + B[1] + B[2])
	if sei == "男" {
		switch S {
		case 3, 4:
			S = 1
		case 5, 6:
			S = 2
		case 7, 8:
			S = 3
		case 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3:
			S = 1
		case 4, 5:
			S = 2
		case 6, 7:
			S = 3
		case 8, 9, 10:
			S = 4
		case 11, 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten17(sei string, c2, c5, c8 int) int {

	B := setBlank3(c2, c5, c8)

	//S := 15 - (c2 + c5 + c8)
	S := 15 - (B[0] + B[1] + B[2])
	if sei == "男" {
		switch S {
		case 3, 4, 5:
			S = 1
		case 6, 7:
			S = 2
		case 8, 9:
			S = 3
		case 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3, 4, 5:
			S = 1
		case 6, 7:
			S = 2
		case 8, 9:
			S = 3
		case 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten18(sei string, c3, c6, c9 int) int {

	B := setBlank3(c3, c6, c9)

	//S := 15 - (c3 + c6 + c9)
	S := 15 - (B[0] + B[1] + B[2])
	if sei == "男" {
		switch S {
		case 3, 4, 5, 6:
			S = 1
		case 7, 8:
			S = 2
		case 9:
			S = 3
		case 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 3, 4, 5, 6:
			S = 1
		case 7, 8:
			S = 2
		case 9:
			S = 3
		case 10, 11:
			S = 4
		case 12:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func setSoten19(sei string, d1, d2 int) int {

	D := setBlank2(d1, d2)

	//S := 10 - (d1 + d2)
	S := 10 - (D[0] + D[1])
	if sei == "男" {
		switch S {
		case 2, 3:
			S = 1
		case 4:
			S = 2
		case 5, 6:
			S = 3
		case 7:
			S = 4
		case 8:
			S = 5
		default:
			S = 0
		}
	} else if sei == "女" {
		switch S {
		case 2, 3:
			S = 1
		case 4:
			S = 2
		case 5, 6:
			S = 3
		case 7:
			S = 4
		case 8:
			S = 5
		default:
			S = 0
		}
	} else {
		S = 0
	}
	return S
}

func intToString(i []int) []string {
	s := make([]string, len(i))
	for n := range i {
		s[n] = strconv.Itoa(i[n])
	}
	return s
}
func setAbvalue(S string) int {
	var R int

	switch S {
	case "そうだ":
		R = 4
	case "まあそうだ":
		R = 3
	case "ややちがう":
		R = 2
	case "ちがう":
		R = 1
	default:
		R = 0
	}
	return R
}

func setCbvalue(S string) int {
	var R int

	switch S {
	case "非常に":
		R = 4
	case "かなり":
		R = 3
	case "多少":
		R = 2
	case "全くない":
		R = 1
	default:
		R = 0
	}
	return R
}

func setBad(S int) int {
	var R int

	if (S == 1) || (S == 2) {
		R = 1
	} else {
		R = 0
	}
	return R
}

func setGood(S int) int {
	var R int

	if (S == 4) || (S == 5) {
		R = 1
	} else {
		R = 0
	}
	return R
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}

func setBlank2(v1, v2 int) []int {

	a := make([]int, 2)
	a[0] = v1
	a[1] = v2

	//未回答が１個あれば、0埋め
	if v1 == 0 {
		v2 = 0
	} else {
		if v2 == 0 {
			v1 = 0
		}
	}

	a[0] = v1
	a[1] = v2

	return a

}

func setBlank3(v1, v2, v3 int) []int {

	a := make([]int, 3)
	a[0] = v1
	a[1] = v2
	a[2] = v3

	cnt := 0
	total := 0
	for _, v := range a {
		if v != 0 {
			cnt++
		}
		total = total + v
	}

	//回答は2/3以上か
	blankV := 0
	if (cnt >= 2) && (cnt != 3) {
		blankV = int(round(float64(total) / float64(cnt)))
		for i := range a {
			if a[i] == 0 {
				a[i] = blankV
			}
		}
	} else {
		//回答は1/3以下なら0埋め
		if cnt != 3 {
			for i := range a {
				a[i] = 0
			}
		}
	}

	return a

}

func setBlank33(v1, v2, v3 int) []int {

	a := make([]int, 3)
	a[0] = v1
	a[1] = v2
	a[2] = v3

	cnt := 0
	total := 0
	for _, v := range a {
		if v != 0 {
			cnt++
		}
		total = total + v
	}

	//回答は2/3以上か
	blankV := 0
	if (cnt >= 2) && (cnt != 3) {
		blankV = int(round(float64(total) / float64(cnt)))
		for i := range a {
			if a[i] == 0 {
				if i == 2 {
					a[i] = 5 - blankV
				} else {
					a[i] = blankV
				}
			}
		}
	} else {
		//回答は1/3以下なら0埋め
		if cnt != 3 {
			for i := range a {
				a[i] = 0
			}
		}

	}

	return a

}

func setBlank6(v1, v2, v3, v4, v5, v6 int) []int {

	a := make([]int, 6)
	a[0] = v1
	a[1] = v2
	a[2] = v3
	a[3] = v4
	a[4] = v5
	a[5] = v6

	cnt := 0
	total := 0
	for _, v := range a {
		if v != 0 {
			cnt++
		}
		total = total + v
	}

	//回答は2/3以上か
	blankV := 0
	if (cnt >= 4) && (cnt != 6) {
		blankV = int(round(float64(total) / float64(cnt)))
		for i := range a {
			if a[i] == 0 {
				a[i] = blankV
			}
		}
	} else {
		//回答は1/3以下なら0埋め
		if cnt != 6 {
			for i := range a {
				a[i] = 0
			}
		}

	}

	return a

}

func setBlank11(v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11 int) []int {

	a := make([]int, 11)
	a[0] = v1
	a[1] = v2
	a[2] = v3
	a[3] = v4
	a[4] = v5
	a[5] = v6
	a[6] = v7
	a[7] = v8
	a[8] = v9
	a[9] = v10
	a[10] = v11

	cnt := 0
	total := 0
	for _, v := range a {
		if v != 0 {
			cnt++
		}
		total = total + v
	}

	//回答は2/3以上か
	blankV := 0
	if (cnt >= 7) && (cnt != 11) {
		blankV = int(round(float64(total) / float64(cnt)))
		for i := range a {
			if a[i] == 0 {
				a[i] = blankV
			}
		}
	} else {
		//回答は1/3以下なら0埋め
		if cnt != 11 {
			for i := range a {
				a[i] = 0
			}
		}

	}

	return a

}
