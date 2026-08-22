<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <div>
        <h1>{{ $t('volunteer.title') }}</h1>
        <p class="subtitle">{{ $t('volunteer.subtitle') }}</p>
      </div>
      <div class="header-badges">
        <span class="role-badge">BÉNÉVOLE</span>
      </div>
    </div>

    <div class="tabs">
      <button :class="{ active: tab === 'rounds' }" @click="tab = 'rounds'">
        {{ $t('volunteer.tabs.rounds') }} ({{ rounds.length }})
      </button>
      <button :class="{ active: tab === 'services' }" @click="tab = 'services'">
        {{ $t('volunteer.tabs.services') }} ({{ myServices.length }})
      </button>
      <button :class="{ active: tab === 'missions' }" @click="tab = 'missions'">
        {{ $t('volunteer.tabs.missions') }} ({{ missions.length }})
      </button>
      <button :class="{ active: tab === 'profile' }" @click="tab = 'profile'">
        {{ $t('volunteer.tabs.profile') }}
      </button>
    </div>

    <div v-if="msg" :class="'alert ' + (msgType === 'ok' ? 'alert-success' : 'alert-error')">
      {{ msg }}
    </div>

    <div v-if="tab === 'rounds'">
      <h2 class="section-title">{{ $t('volunteer.roundsTitle') }}</h2>
      <div v-if="rounds.length === 0" class="empty-state">
        <p>{{ $t('volunteer.noRounds') }}</p>
      </div>
      <div v-else class="rounds-grid">
        <div v-for="round in rounds" :key="round.ID" class="round-box">
          <div class="round-header-row">
            <div>
              <h3 class="round-title">Tournée du {{ formatDate(round.date) }}</h3>
              <p class="round-notes" v-if="round.notes">Instructions : {{ round.notes }}</p>
            </div>
            <div class="round-actions">
              <span :class="'badge status-' + round.status">{{ round.status }}</span>
              <button @click="downloadPDF(round.ID)" class="btn btn-secondary btn-sm">
                {{ $t('common.downloadPdf') }}
              </button>
            </div>
          </div>

          <div class="deliveries-wrapper">
            <h4 class="sub-heading">Livraisons attribuées dans cette tournée ({{ round.deliveries ? round.deliveries.length : 0 }}) :</h4>
            <div v-if="!round.deliveries || round.deliveries.length === 0" class="empty-sub">
              Aucune commande client dans cette tournée.
            </div>
            <div class="table-container" v-else>
              <table>
                <thead>
                  <tr>
                    <th>Destinataire</th>
                    <th>Adresse</th>
                    <th>Statut Livraison</th>
                    <th>Action</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="del in round.deliveries" :key="del.ID">
                    <td><strong>{{ del.recipient_name }}</strong></td>
                    <td>{{ del.recipient_address }}</td>
                    <td>
                      <span :class="'badge status-' + del.status">{{ del.status }}</span>
                    </td>
                    <td>
                      <select v-model="del.status" @change="updateDeliveryStatus(del)" class="select-sm">
                        <option value="pending">pending (En attente)</option>
                        <option value="in_transit">in_transit (En cours)</option>
                        <option value="delivered">delivered (Livrée)</option>
                      </select>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'services'">
      <h2 class="section-title">Proposer une nouvelle formation / atelier thématique</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group">
            <label>Titre de la formation / atelier</label>
            <input v-model="newService.title" placeholder="Ex: Atelier Cuisine Anti-Gaspi, Astuces Zéro Déchet..." />
          </div>
          <div class="form-group">
            <label>Catégorie</label>
            <input v-model="newService.category" placeholder="Cuisine, Sensibilisation, Compostage..." />
          </div>
        </div>
        <div class="form-group">
          <label>Description du contenu</label>
          <textarea v-model="newService.description" rows="3" placeholder="Présentez les objectifs de votre atelier..."></textarea>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Nombre de places max</label>
            <input v-model.number="newService.max_participants" type="number" placeholder="Ex: 15" />
          </div>
          <div class="form-group">
            <label>Date et heure</label>
            <input v-model="newService.date" type="datetime-local" />
          </div>
          <div class="form-group">
            <label>Lieu / Adresse</label>
            <input v-model="newService.location" placeholder="Centre Social, Salle municipale, En ligne..." />
          </div>
        </div>
        <button class="btn btn-primary" @click="createService">
          {{ $t('volunteer.createServiceBtn') }}
        </button>
      </div>

      <h2 class="section-title" style="margin-top: 2rem;">{{ $t('volunteer.servicesTitle') }}</h2>
      <div v-if="myServices.length === 0" class="empty-state">
        <p>{{ $t('volunteer.noServices') }}</p>
      </div>
      <div v-else class="services-list">
        <div v-for="srv in myServices" :key="srv.ID" class="service-box">
          <div class="service-top">
            <div>
              <h3>{{ srv.title }}</h3>
              <p class="service-meta">
                {{ formatDate(srv.date) }} |  {{ srv.location }} |  {{ srv.category }}
              </p>
              <p class="service-desc">{{ srv.description }}</p>
            </div>
            <div class="service-badge-col">
              <span class="badge badge-green">{{ srv.registrations ? srv.registrations.length : 0 }} / {{ srv.max_participants }} inscrits</span>
              <button @click="exportServiceCSV(srv)" class="btn btn-secondary btn-sm" style="margin-top:0.5rem;">
                {{ $t('volunteer.exportAttendees') }}
              </button>
            </div>
          </div>

          <div class="attendees-wrapper" v-if="srv.registrations && srv.registrations.length > 0">
            <h4>{{ $t('volunteer.registeredAttendees') }} :</h4>
            <table class="sub-table">
              <thead>
                <tr>
                  <th>Nom / Prénom</th>
                  <th>Email</th>
                  <th>Date d'inscription</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="reg in srv.registrations" :key="reg.ID">
                  <td>{{ reg.client ? (reg.client.firstname + ' ' + reg.client.lastname) : ('Client #' + reg.client_id) }}</td>
                  <td>{{ reg.client ? reg.client.email : '-' }}</td>
                  <td>{{ formatDate(reg.CreatedAt) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'missions'">
      <h2 class="section-title">{{ $t('volunteer.missionsTitle') }}</h2>
      <div v-if="missions.length === 0" class="empty-state">
        <p>{{ $t('volunteer.noMissions') }}</p>
      </div>
      <div class="table-container" v-else>
        <table>
          <thead>
            <tr>
              <th>Commerçant</th>
              <th>Adresse de collecte</th>
              <th>Date de passage prévue</th>
              <th>Statut</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="m in missions" :key="m.ID">
              <td><strong>{{ m.merchant ? m.merchant.company_name : ('Commerçant #' + m.merchant_id) }}</strong></td>
              <td>{{ m.merchant ? m.merchant.company_address : 'N/A' }}</td>
              <td>{{ formatDate(m.pickup_date) }}</td>
              <td><span :class="'badge status-' + m.status">{{ m.status }}</span></td>
              <td>
                <button v-if="m.status !== 'completed'" @click="updateMissionStatus(m.ID, 'completed')" class="btn btn-primary btn-sm">
                  Valider la collecte
                </button>
                <span v-else class="text-success">Collecte effectuée</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="tab === 'profile'">
      <h2 class="section-title">{{ $t('volunteer.profileTitle') }}</h2>
      <div class="form-section max-w-lg">
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('auth.firstname') }}</label>
            <input v-model="editProfile.firstname" />
          </div>
          <div class="form-group">
            <label>{{ $t('auth.lastname') }}</label>
            <input v-model="editProfile.lastname" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('auth.phone') }}</label>
            <input v-model="editProfile.phone" placeholder="06 12 34 56 78" />
          </div>
          <div class="form-group">
            <label>{{ $t('auth.address') }}</label>
            <input v-model="editProfile.address" placeholder="12 Rue de la Paix, Paris" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('dashboards.zoneArea') }}</label>
            <input v-model="editProfile.zone_area" placeholder="Paris, Lyon..." />
          </div>
          <div class="form-group">
            <label>{{ $t('dashboards.availability') }}</label>
            <input v-model="editProfile.availability" placeholder="Matin, Soir, Week-end..." />
          </div>
        </div>
        <div class="form-group margin-v">
          <label style="display:flex; align-items:center; gap:0.5rem; cursor:pointer;">
            <input type="checkbox" v-model="editProfile.vehicle" style="width:auto;" />
            <strong>{{ $t('dashboards.hasVehicle') }}</strong>
          </label>
        </div>
        <button class="btn btn-primary" @click="saveProfile">
          {{ $t('client.updateProfile') }}
        </button>
      </div>
    </div>

  </div>
</template>

<script>
import api from '../services/api'

export default {
  name: 'VolunteerDashboard',
  data() {
    return {
      tab: 'rounds',
      msg: '',
      msgType: '',
      rounds: [],
      missions: [],
      myServices: [],
      editProfile: {
        firstname: '',
        lastname: '',
        phone: '',
        address: '',
        zone_area: '',
        availability: '',
        vehicle: false
      },
      newService: {
        title: '',
        description: '',
        category: 'Sensibilisation',
        max_participants: 15,
        date: '',
        location: ''
      }
    }
  },
  mounted() {
    this.loadVolunteerData()
    this.loadMyServices()
    this.loadProfile()
  },
  methods: {
    flash(text, type) {
      this.msg = text
      this.msgType = type
      setTimeout(() => this.msg = '', 4500)
    },
    formatDate(d) {
      if (!d) return ''
      return new Date(d).toLocaleDateString('fr-FR', {
        day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit'
      })
    },
    async loadVolunteerData() {
      try {
        const resRounds = await api.get('/volunteer/rounds')
        this.rounds = resRounds.data.rounds || []
        if (resRounds.data.volunteer) {
          this.editProfile.zone_area = resRounds.data.volunteer.zone_area || ''
          this.editProfile.availability = resRounds.data.volunteer.availability || ''
          this.editProfile.vehicle = !!resRounds.data.volunteer.vehicle
        }
      } catch (e) {
        this.rounds = []
      }

      try {
        const resMissions = await api.get('/volunteer/missions')
        this.missions = resMissions.data.missions || []
      } catch (e) {
        this.missions = []
      }
    },
    async loadMyServices() {
      try {
        const res = await api.get('/volunteer/services')
        this.myServices = res.data.services || []
      } catch (e) {
        this.myServices = []
      }
    },
    async loadProfile() {
      try {
        const res = await api.get('/profile')
        const data = res.data.data || {}
        this.editProfile.firstname = data.firstname || ''
        this.editProfile.lastname = data.lastname || ''
        this.editProfile.phone = data.phone || ''
        this.editProfile.address = data.address || ''
      } catch (e) {}
    },
    async createService() {
      if (!this.newService.title || !this.newService.date || !this.newService.location) {
        this.flash('Veuillez remplir au moins le titre, la date et le lieu', 'err')
        return
      }
      try {
        await api.post('/volunteer/services', {
          ...this.newService,
          date: new Date(this.newService.date).toISOString()
        })
        this.flash('Formation / Atelier créé avec succès !', 'ok')
        this.newService = { title: '', description: '', category: 'Sensibilisation', max_participants: 15, date: '', location: '' }
        this.loadMyServices()
      } catch (e) {
        this.flash(e.response?.data?.error || 'Erreur lors de la création de la formation', 'err')
      }
    },
    async saveProfile() {
      try {
        await api.put('/volunteer/profile', this.editProfile)
        this.flash('Votre profil et vos disponibilités ont été mis à jour !', 'ok')
      } catch (e) {
        this.flash(e.response?.data?.error || 'Erreur lors de la mise à jour', 'err')
      }
    },
    async updateDeliveryStatus(del) {
      try {
        await api.put('/volunteer/deliveries/' + del.ID, { status: del.status })
        this.flash('Statut de livraison mis à jour', 'ok')
      } catch (e) {
        this.flash(e.response?.data?.error || 'Erreur lors de la mise à jour', 'err')
      }
    },
    async updateMissionStatus(missionId, newStatus) {
      try {
        await api.put('/volunteer/missions/' + missionId + '/status', { status: newStatus })
        this.flash('Mission de collecte validée !', 'ok')
        this.loadVolunteerData()
      } catch (e) {
        this.flash(e.response?.data?.error || 'Erreur lors de la mise à jour', 'err')
      }
    },
    async downloadPDF(roundId) {
      try {
        const response = await api.get(`/volunteer/rounds/${roundId}/pdf`, { responseType: 'blob' })
        const blob = new Blob([response.data], { type: 'application/pdf' })
        const link = document.createElement('a')
        link.href = window.URL.createObjectURL(blob)
        link.download = `tournee-distribution-${roundId}.pdf`
        link.click()
        this.flash('PDF de la tournée téléchargé avec succès !', 'ok')
      } catch (e) {
        this.flash('Erreur lors du téléchargement du PDF', 'err')
      }
    },
    exportServiceCSV(service) {
      if (!service.registrations || service.registrations.length === 0) {
        this.flash('Aucun inscrit pour cette formation', 'err')
        return
      }
      let csvContent = "data:text/csv;charset=utf-8,Nom,Email,Date d'inscription\n"
      service.registrations.forEach(reg => {
        const name = reg.client ? (reg.client.firstname + ' ' + reg.client.lastname) : ('Client #' + reg.client_id)
        const email = reg.client ? reg.client.email : ''
        const date = new Date(reg.CreatedAt).toLocaleDateString('fr-FR')
        csvContent += `"${name}","${email}","${date}"\n`
      })
      const encodedUri = encodeURI(csvContent)
      const link = document.createElement('a')
      link.setAttribute('href', encodedUri)
      link.setAttribute('download', `inscrits-formation-${service.ID}.csv`)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }
  }
}
</script>

<style scoped>
.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.dashboard-header h1 {
  margin: 0;
  font-family: 'Outfit', sans-serif;
  font-size: 2rem;
  color: #ffffff;
}

.subtitle {
  color: #94a3b8;
  font-size: 0.9rem;
  margin: 0.3rem 0 0 0;
}

.role-badge {
  background: #3b82f6;
  color: #ffffff;
  padding: 0.4rem 1rem;
  border-radius: 20px;
  font-weight: 700;
  font-size: 0.8rem;
  letter-spacing: 0.05em;
}

.rounds-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.round-box {
  background: #1e293b;
  padding: 1.5rem;
  border-radius: 14px;
  border: 1px solid #334155;
}

.round-header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
  border-bottom: 1px solid #334155;
  padding-bottom: 0.8rem;
}

.round-title {
  margin: 0;
  color: #52b788;
  font-size: 1.15rem;
}

.round-notes {
  color: #cbd5e1;
  font-size: 0.88rem;
  margin: 0.3rem 0 0 0;
}

.round-actions {
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.deliveries-wrapper {
  margin-top: 0.8rem;
}

.sub-heading {
  color: #94a3b8;
  font-size: 0.92rem;
  margin-bottom: 0.6rem;
}

.select-sm {
  background: #0f172a;
  color: #ffffff;
  border: 1px solid #475569;
  padding: 0.3rem 0.6rem;
  border-radius: 6px;
  font-size: 0.82rem;
}

.services-list {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

.service-box {
  background: #1e293b;
  padding: 1.5rem;
  border-radius: 14px;
  border: 1px solid #334155;
}

.service-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.service-meta {
  color: #52b788;
  font-size: 0.85rem;
  font-weight: 600;
  margin: 0.3rem 0;
}

.service-desc {
  color: #cbd5e1;
  font-size: 0.9rem;
  margin-top: 0.4rem;
}

.service-badge-col {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.attendees-wrapper {
  margin-top: 1.2rem;
  padding-top: 1rem;
  border-top: 1px dashed #334155;
}

.attendees-wrapper h4 {
  color: #94a3b8;
  font-size: 0.88rem;
  margin-bottom: 0.6rem;
}

.sub-table {
  width: 100%;
  border-collapse: collapse;
}

.sub-table th, .sub-table td {
  padding: 0.5rem 0.8rem;
  border-bottom: 1px solid #334155;
  font-size: 0.85rem;
  text-align: left;
}

.sub-table th {
  color: #64748b;
}

.badge-green {
  background: #2d6a4f;
  color: white;
}

.margin-v {
  margin: 1rem 0;
}
</style>