package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"os"
	"log"
)

/* deklarasi fn sebagai function type yang nanti nya akan diproses di fungsi hitung */
type fn func(float64, float64) float64


/*  Fungsi Main setting koneksi dan handler*/
func main() {
	mux := http.NewServeMux()

	// Define route handlers
	mux.HandleFunc("/", perhitungan)
	

	// baca port 80 dari IP mana aja yang request ke server
	http.ListenAndServe(":80", mux)
}

  


/* fungsi perhitungan mengambil request http dan memanggil fungsi matematika yang sesuai permintaan */
func perhitungan(w http.ResponseWriter, r *http.Request) {
	operasi, bilangan, err := getBilangan(r)
	t := 0.0 // init awal nilai t float desimal
	if err == nil {
		switch operasi {
		case "tambah":
			t = hitung(pertambahan, bilangan)
		case "kurang":
			t = hitung(pengurangan, bilangan)
		case "kali":
			t = hitung(perkalian, bilangan)
		case "bagi":
			t = hitung(pembagian, bilangan)
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Eror Bro"))
			return
		}
		fmt.Fprintf(w, "%f", t) 
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Oppss Salah"))
		return
	}

}

/* getBilangan mengambil dari request http mengoper datanya ke operator dan bilangan. */
func getBilangan(r *http.Request) (operator string, bilangan []float64, err error) {
	expr := strings.Split(r.URL.Path, "/") //pemisahan url dengan tanda /
	operator = expr[1] //slice potongan pertama
	expr = expr[2:len(expr)] // potongan kedua dan seterusnya untuk angka angka yang diinputkan
	for _, e := range expr { 
		v, err := strconv.ParseFloat(strings.TrimSpace(e), 64)
		if err == nil {
			bilangan = append(bilangan, v)
		} else {
			return "", nil, errors.New("Salah Inputan!!")
		}
	}
	return operator, bilangan, nil
}

/* fungsi hitung kalkulator dimana jumlah bilangan akan dimasukan dan dihitung hingga bilangan terakhir diproses dan hasil akhirnya dikembalikan ke total  */
func hitung(operasi fn, bilangan []float64) float64 {
	total := 0.0
	for i, e := range bilangan { 
		if i == 0 {
			total = e
		} else {
			total = operasi(total, e)
		}
	}
	return total
}

/* fungsi operator pertambahan */
func pertambahan(x float64, y float64) float64 {
	return x + y
}

/* fungsi operator pengurangan */
func pengurangan(x float64, y float64) float64 {
	return x - y
}

/* fungsi operator pembagian */
func pembagian(x float64, y float64) float64 {
	return x / y
}

/* fungsi operator perkalian . */
func perkalian(x float64, y float64) float64 {
	return x * y
}

