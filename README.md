# NO MORE WASTE — Plateforme de Gestion de Lutte contre le Gaspillage

Développée dans le cadre de la Mission 1 (PA-Rattrapage). 

Cette application Web complète associe une **API Backend haute performance en Go (Gin/GORM/SQLite)** et une **interface Frontend moderne et réactive en Vue 3 (Vite + vue-i18n)**, conteneurisées avec **Docker & Nginx**.

---

## Lancement Rapide avec Docker

La méthode recommandée pour installer et démarrer l'intégralité du projet en une seule commande :

```bash
# 1. Cloner le projet et se placer à la racine
cd no-more-waste

# 2. Démarrer les conteneurs (Backend Go + Frontend Nginx)
docker compose up --build -d
```

L'application est immédiatement accessible sur :
- **Application Web (Frontend + Nginx)** : [http://localhost](http://localhost) (ou [http://localhost:5173](http://localhost:5173))
- **API Backend Go (Healthcheck)** : [http://localhost:8080/ping](http://localhost:8080/ping)

---

## Lancement Manuel (Sans Docker)

### 1. Démarrer le Backend Go
```bash
cd backend
go run cmd/server/main.go
# L'API démarre sur http://localhost:8080
```

### 2. Démarrer le Frontend Vue 3
```bash
cd Frontend
npm install
npm run dev
# L'interface démarre sur http://localhost:5173
```

---

## Identifiants de Test pour la Soutenance

| Rôle | Email | Mot de passe | Description / Fonctionnalités |
| :--- | :--- | :--- | :--- |
| **Admin** | `admin@nomorewaste.fr` | `admin123` | Back-office complet, gestion des utilisateurs, filtres par rôle, staff, bénévoles, tournées PDF. |
| **Staff** | `staff@nomorewaste.fr` | `staff123` | Gestion logistique, validation des commerçants, création des tournées, export Excel/CSV. |
| **Merchant** | `jean@boulangerie.fr` | `merchant123` | Espace commerçant, ajout de lots invendus, génération automatique de code-barres `NMW-...`. |
| **Client** | `alice@example.com` | `client123` | Front-office bénéficiaires, consultation des produits, inscription aux services, profil. |

---

## Couverture des Exigences de la Mission 1

1. **Adhésions Commerçants & Renouvellement** :
   - Inscription et validation par le staff.
   - Routine automatique en arrière-plan (`renewal_service.go`) qui vérifie quotidiennement les expirations à J-7 et envoie un rappel par e-mail.

2. **Système de Collectes** :
   - Organisation des collectes d'invendus avec affectation des bénévoles et véhicules.

3. **Stockage & Code-barres** :
   - Code-barres unique attribué à chaque lot (`NMW-YYYYMMDD...`).
   - Recherche rapide textuelle pour les clients et interface douchette pour la logistique.

4. **Tournées de Distribution & Génération PDF** :
   - Planification des livraisons (associations / particuliers).
   - Génération dynamique de bordereau de livraison récapitulatif au format **PDF** (`GET /distribution-rounds/:id/pdf`).

5. **Suivi des Bénévoles & Compétences** :
   - Gestion des véhicules, zones de disponibilité et compétences (*Chauffeur, Cuisinier, Bricoleur...*).
   - Sélection contextuelle intelligente dans les formulaires de tournées et services.

6. **Services & Export Excel / CSV** :
   - Publication de services (cours de cuisine anti-gaspillage, aide...) et inscriptions.
   - Route d'exportation de planning au format **CSV / Excel** (`GET /staff/services/export`).

7. **Internationalisation (i18n)** :
   - Multi-langues natif sur le Frontend avec support du **Français (FR)**, **Anglais (EN)** et **Italien (IT)** (Porto, Naples, Dublin).

8. **Infrastructure Nginx & Erreurs HTTP** :
   - Proxy inverse vers l'API Go.
   - Réécriture d'URL SPA (`try_files $uri $uri/ /index.html`).
   - Pages d'erreur personnalisées aux couleurs du projet (`custom_404.html`, `custom_500.html`).

9. **Vérification SIRET en Temps Réel (KYB / API État)** :
   - Interrogation directe de l'API officielle `recherche-entreprises.api.gouv.fr` lors de l'inscription commerçant (`POST /register/merchant`).
   - Rejet automatique des faux numéros SIRET ou des entreprises fermées/en liquidation (`etat_administratif == "F"`).
   - Mode Fallback hors-ligne pour la soutenance si pas de réseau.
