package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
)

func creerUtilisateur(nom string) {
	numeroCompte := rand.Int31()
	compte := creationCompte(nom, numeroCompte)
	date := dateAujourdhui()
	entete := creationEntete(date, 200)

	var fichierBanque *os.File
	var banque Banque

	fichierBanque, _ = os.OpenFile("banque/banque/banque.dat", os.O_RDWR, 0644)

	defer func(fichierBanque *os.File) {
		err := fichierBanque.Close()
		if err != nil {

		}
	}(fichierBanque)

	if utilisateurEnregistre(fichierBanque, nom, numeroCompte) {
		return
	}

	creationFichierCompte(entete, numeroCompte)

	FReadBanque(fichierBanque, &banque, 0, 0)
	banque.NombreClients += 1

	FWrite(fichierBanque, banque, 0, 0)
	FWrite(fichierBanque, compte, 0, 2)
}

func utilisateurEnregistre(fichierBanque *os.File, nom string, numeroCompte int32) bool {
	var compte Compte
	var banque Banque

	FReadBanque(fichierBanque, &banque, 0, 0)
	nombreClients := int(banque.NombreClients)

	for i := 0; i < nombreClients; i++ {
		FReadCompte(fichierBanque, &compte, lenBanque+int64(i)*lenCompte, 0)
		if compte.Numero == numeroCompte || comparaisonString(string(compte.Nom[:]), nom) {
			fmt.Printf("Le client %s dispose déjà d'un compte !\n", nom)
			return true
		}
	}

	return false
}

func compteDe(nom string) int32 {
	var fichierBanque *os.File
	var banque Banque
	var compte Compte

	fichierBanque, _ = os.OpenFile("banque/banque/banque.dat", os.O_RDWR, 0644)

	defer func(fichierBanque *os.File) {
		err := fichierBanque.Close()
		if err != nil {

		}
	}(fichierBanque)

	FReadBanque(fichierBanque, &banque, 0, io.SeekStart)
	nombreClients := int(banque.NombreClients)

	for i := 0; i < nombreClients; i++ {
		FReadCompte(fichierBanque, &compte, lenBanque+int64(i)*lenCompte, 0)
		if comparaisonString(string(compte.Nom[:]), nom) {
			return compte.Numero
		}
	}

	return 0
}

func compteNumero(numeroCompte int32) string {
	var fichierBanque *os.File
	var banque Banque
	var compte Compte

	fichierBanque, _ = os.OpenFile("banque/banque/banque.dat", os.O_RDWR, 0644)

	defer func(fichierBanque *os.File) {
		err := fichierBanque.Close()
		if err != nil {

		}
	}(fichierBanque)

	FReadBanque(fichierBanque, &banque, 0, io.SeekStart)
	nombreClients := int(banque.NombreClients)

	for i := 0; i < nombreClients; i++ {
		FReadCompte(fichierBanque, &compte, lenBanque+int64(i)*lenCompte, 0)
		if numeroCompte == compte.Numero {
			return string(compte.Nom[:])
		}
	}

	return ""
}

func soldeDe(nom string, date Date) float32 {
	var fichier *os.File
	var entete Entete

	numeroCompte := compteDe(nom)
	ouvrirFichierCompte(&fichier, int(numeroCompte))

	defer func(fichier *os.File) {
		err := fichier.Close()
		if err != nil {

		}
	}(fichier)

	misaAJourSolde(fichier, date)

	FReadEntete(fichier, &entete, 0, io.SeekStart)

	return entete.Solde
}

func virementCompteACompte(numeroCompte1 int32, numeroCompte2 int32, date Date, montant float32, label string) {
	var fichier1, fichier2 *os.File

	ouvrirFichierCompte(&fichier1, int(numeroCompte1))
	ouvrirFichierCompte(&fichier2, int(numeroCompte2))

	defer func(fichier1 *os.File) {
		err := fichier1.Close()
		if err != nil {

		}
	}(fichier1)
	defer func(fichier2 *os.File) {
		err := fichier2.Close()
		if err != nil {

		}
	}(fichier2)

	nom1 := compteNumero(numeroCompte1)
	nom2 := compteNumero(numeroCompte2)

	transaction1 := creationTransaction(date, -montant, label, nom2)
	transaction2 := creationTransaction(date, montant, label, nom1)

	ajouterTransaction(fichier1, transaction1)
	ajouterTransaction(fichier2, transaction2)
}

func virementPersonneAPersonne(nom1 string, nom2 string, date Date, montant float32, label string) {
	numeroCompte1 := compteDe(nom1)
	numeroCompte2 := compteDe(nom2)

	virementCompteACompte(numeroCompte1, numeroCompte2, date, montant, label)
}

func imprimerReleve(nom string, date Date) {
	var fichierCompte, fichierReleve *os.File
	var entete Entete
	var transaction Transaction

	numeroCompte := compteDe(nom)

	ouvrirFichierCompte(&fichierCompte, int(numeroCompte))
	ouvrirFichierReleve(&fichierReleve, int(numeroCompte))

	defer func(fichierCompte *os.File) {
		err := fichierCompte.Close()
		if err != nil {

		}
	}(fichierCompte)
	defer func(fichierReleve *os.File) {
		err := fichierReleve.Close()
		if err != nil {

		}
	}(fichierReleve)

	misaAJourSolde(fichierCompte, date)

	FReadEntete(fichierCompte, &entete, 0, 0)
	SWriteEntete(fichierReleve, entete)

	for i := int32(0); i < entete.NombreTransaction; i++ {
		FReadTransaction(fichierCompte, &transaction, lenEntete+int64(i)*(lenTransaction-2), 0)
		SWriteTransaction(fichierReleve, transaction)
	}
}

func commandeAjout() {
	var nom string

	fmt.Print("Nom du client : ")
	_, err := fmt.Scanf("%s\n", &nom)
	if err != nil {
		return
	}

	fmt.Print(nom)

	creerUtilisateur(nom)
}

func commandeLister() {
	var fichierBanque *os.File
	var banque Banque
	var compte Compte

	fichierBanque, _ = os.OpenFile("banque/banque/banque.dat", os.O_RDWR, 0644)

	defer func(fichierBanque *os.File) {
		err := fichierBanque.Close()
		if err != nil {

		}
	}(fichierBanque)

	FReadBanque(fichierBanque, &banque, 0, io.SeekStart)
	nombreClients := int(banque.NombreClients)

	for i := 0; i < nombreClients; i++ {
		FReadCompte(fichierBanque, &compte, lenBanque+int64(i)*lenCompte, 0)
		fmt.Printf("Client : %s, Numero de compte : %d\n", bytes.Trim(compte.Nom[:], "\x00"), compte.Numero)
	}
}

func commandeVirement() {
	var emmeteur, receveur, label string
	var montant float32

	fmt.Print("Émetteur : ")
	_, err := fmt.Scanf("%s\n", &emmeteur)
	if err != nil {
		return
	}
	fmt.Print("Receveur  : ")
	_, err = fmt.Scanf("%s\n", &receveur)
	if err != nil {
		return
	}
	fmt.Print("Montant de la transaction : ")
	_, err = fmt.Scanf("%f\n", &montant)
	if err != nil {
		return
	}
	fmt.Print("Descriptif de la transaction : ")
	_, err = fmt.Scanf("%s\n", &label)
	if err != nil {
		return
	}

	virementPersonneAPersonne(emmeteur, receveur, dateAujourdhui(), montant, label)
}

func commandeMiseAJour() {
	var nom string
	var fichier *os.File

	fmt.Print("Votre nom : ")
	_, err := fmt.Scanf("%s\n", &nom)
	if err != nil {
		return
	}

	ouvrirFichierCompte(&fichier, int(compteDe(nom)))

	defer func(fichier *os.File) {
		err := fichier.Close()
		if err != nil {

		}
	}(fichier)

	misaAJourSolde(fichier, dateAujourdhui())

	fmt.Printf("%s disose de %.2f€ sur son compte aujourd'hui !\n", nom, soldeDe(nom, dateAujourdhui()))
}

func commandeReleve() {
	var nom string
	var fichier *os.File

	fmt.Print("Votre nom : ")
	_, err := fmt.Scanf("%s\n", &nom)
	if err != nil {
		return
	}

	ouvrirFichierCompte(&fichier, int(compteDe(nom)))

	defer func(fichier *os.File) {
		err := fichier.Close()
		if err != nil {

		}
	}(fichier)

	imprimerReleve(nom, dateAujourdhui())
}

func menu() {
	if _, err := os.Stat("banque/banque/banque.dat"); errors.Is(err, os.ErrNotExist) {
		fmt.Print("Création de la banque\n")
		creationFichierBanque()
	}

	var commande string
Loop:
	for {
		fmt.Print("\n\nAjouter un nouveau client..............: A\n")
		fmt.Print("Lister tous les comptes de clients.....: L\n")
		fmt.Print("Virement depuis un compte client.......: V\n")
		fmt.Print("Mise à jour du solde d'un client.......: M\n")
		fmt.Print("Relevé d'un compte client..............: R\n")
		fmt.Print("Quitter................................: Q\n")
		fmt.Print("Votre choix : ")
		_, err := fmt.Scanf("%s\n", &commande)
		if err != nil {
			return
		}
		switch commande {
		case "A", "a":
			commandeAjout()
		case "L", "l":
			commandeLister()
		case "V", "v":
			commandeVirement()
		case "M", "m":
			commandeMiseAJour()
		case "R", "r":
			commandeReleve()
		case "Q", "q":
			break Loop
		default:
			fmt.Print("Commande invalide !")
		}
	}
}
