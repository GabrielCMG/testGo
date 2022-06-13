package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
	"unsafe"
)

type Date struct {
	Jour  int32
	Mois  int32
	Annee int32
}

type Compte struct {
	Nom    [20]byte
	Numero int32
}

type Transaction struct {
	Date    Date
	Montant float32
	Label   [50]byte
	Nom     [20]byte
}

type Entete struct {
	Date              Date
	Solde             float32
	NombreTransaction int32
}

type Banque struct {
	NombreClients int32
}

const (
	lenEntete      = int64(unsafe.Sizeof(Entete{}))
	lenTransaction = int64(unsafe.Sizeof(Transaction{}))
	lenCompte      = int64(unsafe.Sizeof(Compte{}))
	lenBanque      = int64(unsafe.Sizeof(Banque{}))
)

func dateAujourdhui() Date {
	var date Date

	now := time.Now()

	date.Jour = int32(now.Day())
	date.Mois = int32(now.Month())
	date.Annee = int32(now.Year())

	return date
}

func comparaisonDates(date1 Date, date2 Date) int {
	if date1.Annee == date2.Annee && date1.Mois == date2.Mois {
		if date1.Jour == date2.Jour {
			return 0
		} else if date1.Jour > date2.Jour {
			return 1
		} else {
			return -1
		}
	}
	if date1.Annee == date2.Annee {
		if date1.Mois > date2.Mois {
			return 1
		} else if date1.Mois < date2.Mois {
			return -1
		}
	}
	if date1.Annee > date2.Annee {
		return 1
	} else if date1.Annee < date2.Annee {
		return -1
	}
	return 0
}

func comparaisonString(string1 string, string2 string) bool {
	var string1B [50]byte
	var string2B [50]byte

	copy(string1B[:], string1)
	copy(string2B[:], string2)

	return bytes.Compare(string1B[:], string2B[:]) == 0
}

func toBytes(o any) []byte {
	var binBuf bytes.Buffer
	err := binary.Write(&binBuf, binary.BigEndian, o)

	if err != nil {
		return nil
	}

	return binBuf.Bytes()
}

func FWrite(fichier *os.File, objet any, offset int64, whence int) {
	_, err := fichier.Seek(offset, whence)
	if err != nil {
		return
	}

	objetB := toBytes(objet)
	_, err = fichier.Write(objetB)
	if err != nil {
		return
	}
}

func SWriteEntete(fichier *os.File, entete Entete) {
	texte := fmt.Sprintf("Argent sur le compte le %d/%d/%d : %.2f€\n",
		entete.Date.Jour, entete.Date.Mois, entete.Date.Annee, entete.Solde)

	_, err := fichier.WriteString(texte)
	if err != nil {
		return
	}
}

func SWriteTransaction(fichier *os.File, transaction Transaction) {
	var texte string

	var nom, label []byte
	label = bytes.Trim(transaction.Label[:], "\x00")
	nom = bytes.Trim(transaction.Nom[:], "\x00")

	if transaction.Montant >= 0 {
		texte = fmt.Sprintf("Virement du %d/%d/%d reçu de la part de %s\n   Montant : %.2f\n   Description : %s\n",
			transaction.Date.Jour, transaction.Date.Mois, transaction.Date.Annee,
			nom, transaction.Montant, label)
	} else {
		texte = fmt.Sprintf("Virement du %d/%d/%d effectué à %s\n   Montant : %.2f\n   Description : %s\n",
			transaction.Date.Jour, transaction.Date.Mois, transaction.Date.Annee,
			nom, transaction.Montant, label)
	}

	fmt.Print(texte)

	_, err := fichier.WriteString(texte)
	if err != nil {
		return
	}
}

func FReadEntete(fichier *os.File, objet *Entete, offset int64, whence int) Entete {
	_, err := fichier.Seek(offset, whence)
	if err != nil {
		return Entete{}
	}

	objetB := make([]byte, lenEntete)

	_, err = fichier.Read(objetB)
	if err != nil {
		return Entete{}
	}

	buffer := bytes.NewBuffer(objetB)
	err = binary.Read(buffer, binary.BigEndian, objet)
	if err != nil {
		return Entete{}
	}

	return *objet
}

func FReadTransaction(fichier *os.File, objet *Transaction, offset int64, whence int) Transaction {
	_, err := fichier.Seek(offset, whence)
	if err != nil {
		return Transaction{}
	}

	objetB := make([]byte, lenTransaction)

	_, err = fichier.Read(objetB)
	if err != nil {
		return Transaction{}
	}

	buffer := bytes.NewBuffer(objetB)
	err = binary.Read(buffer, binary.BigEndian, objet)
	if err != nil {
		return Transaction{}
	}

	return *objet
}

func FReadCompte(fichier *os.File, objet *Compte, offset int64, whence int) {
	_, err := fichier.Seek(offset, whence)
	if err != nil {
		return
	}

	objetB := make([]byte, lenCompte)

	_, err = fichier.Read(objetB)
	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(objetB)
	err = binary.Read(buffer, binary.BigEndian, objet)
	if err != nil {
		return
	}
}

func FReadBanque(fichier *os.File, objet *Banque, offset int64, whence int) {
	_, err := fichier.Seek(offset, whence)
	if err != nil {
		return
	}

	objetB := make([]byte, lenCompte)

	_, err = fichier.Read(objetB)
	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(objetB)
	err = binary.Read(buffer, binary.BigEndian, objet)
	if err != nil {
		return
	}
}
