# 🚀 No More Waste — Plateforme en Ligne & Comptes de Démonstration

## 🌐 Liens d'accès à la Plateforme Déployée

- **Site Web Officiel** : [http://nomorewaste.online](http://nomorewaste.online)
- **Adresse IP Directe (VPS)** : [http://92.113.31.151](http://92.113.31.151)
- **API Backend (Healthcheck)** : [http://nomorewaste.online/ping](http://nomorewaste.online/ping)

---

## 🔑 Mot de passe général par défaut
Tous les comptes ci-dessous sont configurés avec le mot de passe unique :
> **`Password123!`**

---

## 📋 Tableau Récapitulatif des Comptes de Production

### 1. 🛡️ Administrateur System
| Rôle | Email | Mot de passe | Rôle / Accès |
| :--- | :--- | :--- | :--- |
| **Admin** | `admin@nomorewaste.com` | `Password123!` | Gestion globale de la plateforme, gestion des utilisateurs, validation des commerçants, statistiques et supervision. |

---

### 2. 💼 Personnel / Staff (Salariés No More Waste)
| Nom & Prénom | Email | Mot de passe | Poste / Département |
| :--- | :--- | :--- | :--- |
| **Marc Dubois** | `staff@nomorewaste.com` | `Password123!` | Responsable de Collecte (*Logistique & Opérations*) |

---

### 3. 🤝 Bénévoles (Volunteers)
| Nom & Prénom | Email | Mot de passe | Zone d'action | Véhiculé | Compétences associées |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Sophie Martin** | `benevole1@nomorewaste.com` | `Password123!` | Paris-Nord | ✅ Oui | Permis B / Transport, Sensibilisation Zéro-Déchet |
| **Lucas Bernard** | `benevole2@nomorewaste.com` | `Password123!` | Paris-Sud | ❌ Non | Distribution Alimentaire, Animation d'Atelier |

---

### 4. 🛒 Clients / Bénéficiaires
| Nom & Prénom | Email | Mot de passe | Adresse | Téléphone |
| :--- | :--- | :--- | :--- | :--- |
| **Claire Petit** | `client1@nomorewaste.com` | `Password123!` | 14 Rue de Rivoli, 75004 Paris | `0711223344` |
| **Antoine Rousseau** | `client2@nomorewaste.com` | `Password123!` | 5 Avenue de la République, 75011 Paris | `0722334455` |
| **Élodie Moreau** | `client3@nomorewaste.com` | `Password123!` | 42 Rue de Clichy, 75009 Paris | `0733445566` |

---

## 🛠️ Commande pour Réexécuter le Seed (Réinitialiser les données)

Si vous souhaitez réinitialiser ou peupler à nouveau la base de données sur le VPS :

```bash
curl -X POST http://localhost:8080/seed
```
