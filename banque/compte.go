package main

import (
	"os"
	"strconv"
)

func creationEntete(date Date, solde float32) Entete {
	var entete Entete

	entete.Date = date
	entete.Solde = solde
	entete.NombreTransaction = 0

	return entete
}

func creationTransaction(date Date, montant float32, label string, nom string) Transaction {
	var transaction Transaction

	transaction.Date = date
	transaction.Montant = montant
	copy(transaction.Nom[:], nom)
	copy(transaction.Label[:], label)

	return transaction
}

func creationCompte(nom string, numeroCompte int32) Compte {
	var compte Compte

	copy(compte.Nom[:], nom)
	compte.Numero = numeroCompte

	return compte
}

func creationFichierCompte(entete Entete, noCompte int32) {
	var fichier *os.File

	nomFichier := "banque/banque/" + strconv.Itoa(int(noCompte)) + ".dat"
	fichier, _ = os.OpenFile(nomFichier, os.O_CREATE, 0644)

	defer func(fichier *os.File) {
		err := fichier.Close()
		if err != nil {

		}
	}(fichier)

	FWrite(fichier, entete, 0, 0)
}

func creationFichierBanque() {
	var fichier *os.File

	nomFichier := "banque/banque/banque.dat"
	fichier, _ = os.OpenFile(nomFichier, os.O_CREATE, 0644)

	defer func(fichier *os.File) {
		err := fichier.Close()
		if err != nil {

		}
	}(fichier)

	FWrite(fichier, Banque{0}, 0, 0)
}

func ouvrirFichierCompte(fichier **os.File, noCompte int) {
	nomFichier := "banque/banque/" + strconv.Itoa(noCompte) + ".dat"

	*fichier, _ = os.OpenFile(nomFichier, os.O_RDWR, 0644)
}

func ouvrirFichierReleve(fichier **os.File, noCompte int) {
	nomFichier := "banque/relevés/releve_" + strconv.Itoa(noCompte) + ".txt"

	*fichier, _ = os.OpenFile(nomFichier, os.O_RDWR|os.O_CREATE, 0644)
}

func ajouterTransaction(fichier *os.File, transaction Transaction) {
	var entete Entete

	entete = FReadEntete(fichier, &entete, 0, 0)
	entete.NombreTransaction += 1

	FWrite(fichier, entete, 0, 0)
	FWrite(fichier, transaction, 0, 2)
}

func misaAJourSolde(fichier *os.File, date Date) {
	var entete Entete
	var transaction Transaction

	FReadEntete(fichier, &entete, 0, 0)
	nouveauSolde := entete.Solde

	for i := 0; i < int(entete.NombreTransaction); i++ {
		FReadTransaction(fichier, &transaction, lenEntete+int64(i)*(lenTransaction-2), 0)

		if comparaisonDates(entete.Date, transaction.Date) <= 0 &&
			comparaisonDates(transaction.Date, date) < 0 &&
			comparaisonDates(entete.Date, date) < 0 {
			nouveauSolde += transaction.Montant
		}
		if comparaisonDates(date, transaction.Date) <= 0 &&
			comparaisonDates(transaction.Date, entete.Date) < 0 &&
			comparaisonDates(date, entete.Date) < 0 {
			nouveauSolde -= transaction.Montant
		}
	}

	entete.Solde = nouveauSolde
	entete.Date = date
	FWrite(fichier, entete, 0, 0)
}
