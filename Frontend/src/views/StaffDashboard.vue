<template>
  <div class="dashboard">
    <h1>{{ $t('dashboards.staffTitle') }}</h1>

    <div class="tabs">
      <button :class="{ active: tab === 'merchants' }" @click="tab = 'merchants'">{{ $t('dashboards.tabs.merchants') }}</button>
      <button :class="{ active: tab === 'volunteers' }" @click="tab = 'volunteers'">{{ $t('dashboards.tabs.volunteers') }}</button>
      <button :class="{ active: tab === 'skills' }" @click="tab = 'skills'">{{ $t('dashboards.tabs.skills') }}</button>
      <button :class="{ active: tab === 'missions' }" @click="tab = 'missions'">{{ $t('dashboards.tabs.missions') }}</button>
      <button :class="{ active: tab === 'services' }" @click="tab = 'services'">{{ $t('dashboards.tabs.services') }}</button>
      <button :class="{ active: tab === 'rounds' }" @click="tab = 'rounds'">{{ $t('dashboards.tabs.rounds') }}</button>
      <button :class="{ active: tab === 'deliveries' }" @click="tab = 'deliveries'">Paniers Clients Payés</button>
    </div>

    <div v-if="msg" :class="'alert ' + (msgType === 'ok' ? 'alert-success' : 'alert-error')">{{ msg }}</div>

    <div v-if="tab === 'merchants'">
      <h2 class="section-title">Commerçants en attente de validation</h2>
      <div v-if="pendingMerchants.length === 0" class="empty-state">Aucun commerçant en attente</div>
      <div class="table-container" v-else>
        <table>
          <thead><tr><th>Entreprise</th><th>SIRET</th><th>Adresse</th><th>Contact</th><th>Actions</th></tr></thead>
          <tbody>
            <tr v-for="m in pendingMerchants" :key="m.ID">
              <td>{{ m.company_name }}</td>
              <td>{{ m.siret_number }}</td>
              <td>{{ m.company_address }}</td>
              <td>{{ m.user?.firstname }} {{ m.user?.last_name }}</td>
              <td><button class="btn btn-primary btn-sm" @click="approveMerchant(m.ID)">Approuver</button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="tab === 'volunteers'">
      <h2 class="section-title">Inscrire un bénévole</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group"><label>Prénom</label><input v-model="newVolunteer.firstname" /></div>
          <div class="form-group"><label>Nom</label><input v-model="newVolunteer.lastname" /></div>
        </div>
        <div class="form-row">
          <div class="form-group"><label>Email</label><input v-model="newVolunteer.email" type="email" /></div>
          <div class="form-group"><label>Mot de passe</label><input v-model="newVolunteer.password" type="password" /></div>
        </div>
        <div class="form-row">
          <div class="form-group"><label>Zone</label><input v-model="newVolunteer.zone_area" placeholder="Paris, Nantes..." /></div>
          <div class="form-group"><label>Disponibilité</label><input v-model="newVolunteer.availability" placeholder="matin, après-midi..." /></div>
        </div>
        <div class="form-group">
          <label><input type="checkbox" v-model="newVolunteer.vehicle" /> Possède un véhicule</label>
        </div>
        <button class="btn btn-primary" @click="registerVolunteer">Inscrire</button>
      </div>

      <h2 class="section-title" style="margin-top:2rem;">Liste des bénévoles</h2>
      <div v-if="volunteers.length === 0" class="empty-state">Aucun bénévole</div>
      <div class="table-container" v-else>
        <table>
          <thead><tr><th>Nom</th><th>Zone</th><th>Disponibilité</th><th>Véhicule</th><th>Compétences</th><th>Actions</th></tr></thead>
          <tbody>
            <tr v-for="v in volunteers" :key="v.ID">
              <td>{{ v.user?.firstname }} {{ v.user?.last_name }}</td>
              <td>{{ v.zone_area }}</td>
              <td>{{ v.availability }}</td>
              <td>{{ v.vehicle ? 'Oui' : 'Non' }}</td>
              <td><span class="badge badge-blue" v-for="s in v.skills" :key="s.ID" style="margin-right:4px;">{{ s.name }}</span></td>
              <td>
                <select v-model="skillAssign[v.ID]" multiple style="font-size:0.8rem;min-width:120px;">
                  <option v-for="sk in allSkills" :key="sk.ID" :value="sk.ID">{{ sk.name }}</option>
                </select>
                <button class="btn btn-primary btn-sm" style="margin-left:4px;" @click="assignSkills(v.ID)">Assigner</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="tab === 'skills'">
      <h2 class="section-title">Gérer les compétences</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group"><label>Nom</label><input v-model="newSkill.name" placeholder="Chauffeur, Cuisinier..." /></div>
          <div class="form-group"><label>Catégorie</label><input v-model="newSkill.category" placeholder="Transport, Cuisine..." /></div>
        </div>
        <button class="btn btn-primary" @click="createSkill">Ajouter</button>
      </div>
      <div class="card-grid" style="margin-top:1rem;">
        <div class="card" v-for="s in allSkills" :key="s.ID">
          <h3>{{ s.name }}</h3>
          <p>Catégorie : {{ s.category }}</p>
          <div class="btn-group"><button class="btn btn-danger btn-sm" @click="deleteSkill(s.ID)">Supprimer</button></div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'missions'">
      <h2 class="section-title">Créer une mission de collecte</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group">
            <label>ID Commerçant</label>
            <input v-model.number="newMission.merchant_id" type="number" placeholder="ID Commerçant (ex: 1)" />
          </div>
          <div class="form-group">
            <label>Bénévole (Chauffeur / Collecteur)</label>
            <select v-model.number="newMission.volunteer_id">
              <option :value="null">-- Sélectionner un bénévole --</option>
              <option v-for="v in volunteers" :key="v.ID" :value="v.ID">
                {{ v.user?.firstname }} {{ v.user?.last_name }} ({{ v.vehicle ? 'Chauffeur avec Véhicule' : 'Sans Véhicule' }} — {{ v.skills?.map(s => s.name).join(', ') || 'Bénévole' }})
              </option>
            </select>
          </div>
        </div>
        <div class="form-row">
          <div class="form-group"><label>ID Produit</label><input v-model.number="newMission.product_id" type="number" /></div>
          <div class="form-group"><label>Date de collecte</label><input v-model="newMission.pickup_date" type="date" /></div>
        </div>
        <button class="btn btn-primary" @click="createMission">Créer la mission</button>
      </div>
    </div>

    <div v-if="tab === 'services'">
      <h2 class="section-title">Créer un service</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group"><label>Titre</label><input v-model="newService.title" /></div>
          <div class="form-group"><label>Catégorie</label><input v-model="newService.category" placeholder="cuisine, bricolage..." /></div>
        </div>
        <div class="form-group"><label>Description</label><textarea v-model="newService.description" rows="2"></textarea></div>
        <div class="form-row">
          <div class="form-group"><label>Lieu</label><input v-model="newService.location" /></div>
          <div class="form-group"><label>Places max</label><input v-model.number="newService.max_participants" type="number" /></div>
        </div>
        <div class="form-row">
          <div class="form-group"><label>Date</label><input v-model="newService.date" type="date" /></div>
          <div class="form-group">
            <label>Bénévole Animateur</label>
            <select v-model.number="newService.volunteer_id">
              <option :value="null">-- Sélectionner un bénévole --</option>
              <option v-for="v in volunteers" :key="v.ID" :value="v.ID">
                {{ v.user?.firstname }} {{ v.user?.last_name }} (Spécialité : {{ v.skills?.map(s => s.name).join(', ') || 'Polyvalent' }})
              </option>
            </select>
          </div>
        </div>
        <button class="btn btn-primary" @click="createService">Créer</button>
      </div>

      <div style="display: flex; justify-content: space-between; align-items: center; margin-top:2rem;">
        <h2 class="section-title">Tous les services (Plannings)</h2>
        <button class="btn btn-secondary" @click="exportServicesCSV">
          Exporter en Excel (CSV)
        </button>
      </div>
      <div v-if="servicesList.length === 0" class="empty-state">Aucun service</div>
      <div class="card-grid">
        <div class="card" v-for="s in servicesList" :key="s.ID">
          <h3>{{ s.title }}</h3>
          <p>{{ s.description }}</p>
          <p><span class="badge badge-blue">{{ s.category }}</span> <span class="badge" :class="s.status === 'open' ? 'badge-green' : s.status === 'full' ? 'badge-orange' : 'badge-gray'">{{ s.status }}</span></p>
          <p><strong>Date :</strong> {{ formatDate(s.date) }} | <strong>Lieu :</strong> {{ s.location }}</p>
          <p><strong>Inscrits :</strong> {{ s.registrations?.length || 0 }} / {{ s.max_participants }}</p>
          <div class="btn-group">
            <button class="btn btn-danger btn-sm" @click="deleteService(s.ID)">Supprimer</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'rounds'">
      <h2 class="section-title">Créer une tournée</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group"><label>Date</label><input v-model="newRound.date" type="date" /></div>
          <div class="form-group">
            <label>Bénévole Chauffeur (Transport / Livraison)</label>
            <select v-model.number="newRound.volunteer_id">
              <option :value="null">-- Sélectionner un chauffeur --</option>
              <option v-for="v in volunteers" :key="v.ID" :value="v.ID">
                {{ v.user?.firstname }} {{ v.user?.last_name }} {{ v.vehicle ? '(Véhicule OK)' : '(Sans véhicule)' }} — {{ v.skills?.map(s => s.name).join(', ') || 'Chauffeur' }}
              </option>
            </select>
          </div>
        </div>
        <div class="form-group"><label>Notes</label><input v-model="newRound.notes" placeholder="Instructions spéciales..." /></div>
        <button class="btn btn-primary" @click="createRound">Créer la tournée</button>
      </div>

      <h2 class="section-title" style="margin-top:2rem;">Tournées existantes</h2>
      <div v-if="rounds.length === 0" class="empty-state">Aucune tournée</div>
      <div class="card" v-for="r in rounds" :key="r.ID" style="margin-bottom:1.5rem;">
        <div style="display:flex;justify-content:space-between;align-items:center;">
          <h3>Tournée #{{ r.ID }} — {{ formatDate(r.date) }}</h3>
          <div class="btn-group">
            <span class="badge" :class="r.status === 'planned' ? 'badge-blue' : r.status === 'in_progress' ? 'badge-orange' : 'badge-green'">{{ r.status }}</span>
            <button class="btn btn-secondary btn-sm" @click="downloadPDF(r.ID)">Télécharger PDF</button>
            <button class="btn btn-danger btn-sm" @click="deleteRound(r.ID)">Supprimer</button>
          </div>
        </div>
        <p><strong>Bénévole :</strong> {{ r.volunteer?.user?.firstname }} {{ r.volunteer?.user?.last_name }}</p>
        <p v-if="r.notes"><strong>Notes :</strong> {{ r.notes }}</p>

        <div style="margin-top:0.8rem;">
          <strong>Statut :</strong>
          <select v-model="roundStatus[r.ID]" style="margin-left:4px;padding:4px;font-size:0.8rem;">
            <option value="planned">planned</option>
            <option value="in_progress">in_progress</option>
            <option value="completed">completed</option>
          </select>
          <button class="btn btn-primary btn-sm" style="margin-left:4px;" @click="updateRoundStatus(r.ID)">Modifier</button>
        </div>

        <div style="margin-top:1rem;">
          <strong>Livraisons ({{ r.deliveries?.length || 0 }}) :</strong>
          <table v-if="r.deliveries?.length" style="margin-top:0.5rem;">
            <thead><tr><th>Destinataire</th><th>Type</th><th>Adresse</th><th>Produit</th><th>Qté</th><th>Statut</th></tr></thead>
            <tbody>
              <tr v-for="d in r.deliveries" :key="d.ID">
                <td>{{ d.recipient_name }}</td>
                <td><span class="badge" :class="d.recipient_type === 'association' ? 'badge-blue' : 'badge-gray'">{{ d.recipient_type }}</span></td>
                <td>{{ d.recipient_address }}</td>
                <td>{{ d.product?.title }}</td>
                <td>{{ d.quantity }}</td>
                <td><span class="badge" :class="d.status === 'delivered' ? 'badge-green' : 'badge-orange'">{{ d.status }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="form-section" style="margin-top:1rem;background:#f8f9fa;">
          <strong>Ajouter une livraison :</strong>
          <div class="form-row" style="margin-top:0.5rem;">
            <div class="form-group"><label>Destinataire</label><input v-model="newDelivery[r.ID + '_name']" /></div>
            <div class="form-group"><label>Adresse</label><input v-model="newDelivery[r.ID + '_addr']" /></div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Type</label>
              <select v-model="newDelivery[r.ID + '_type']"><option value="association">Association</option><option value="particulier">Particulier</option></select>
            </div>
            <div class="form-group"><label>ID Produit</label><input v-model.number="newDelivery[r.ID + '_pid']" type="number" /></div>
          </div>
          <div class="form-group"><label>Quantité</label><input v-model.number="newDelivery[r.ID + '_qty']" type="number" /></div>
          <button class="btn btn-primary btn-sm" @click="addDelivery(r.ID)">Ajouter</button>
        </div>
      </div>
    </div>

    <div v-if="tab === 'deliveries'">
      <h2 class="section-title">Paniers / Commandes Clients Validés (Paiement Stripe OK)</h2>
      <p style="color: #64748b; margin-bottom: 1rem;">
        Consultez les paniers payés par les clients et assignez-les aux tournées des bénévoles selon leurs zones et disponibilités.
      </p>

      <div v-if="deliveries.length === 0" class="empty-state">Aucun panier client en attente de livraison</div>
      <div class="table-container" v-else>
        <table>
          <thead>
            <tr>
              <th>Client / Destinataire</th>
              <th>Adresse de Livraison</th>
              <th>Produit Commandé</th>
              <th>Quantité</th>
              <th>Statut Livraison</th>
              <th>Tournée Assignée</th>
              <th>Action Logistique</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="d in deliveries" :key="d.ID">
              <td><strong>{{ d.recipient_name }}</strong></td>
              <td>{{ d.recipient_address }}</td>
              <td>{{ d.product?.title || 'Produit #' + d.product_id }}</td>
              <td>{{ d.quantity }}</td>
              <td>
                <span class="badge" :class="d.status === 'delivered' ? 'badge-green' : 'badge-orange'">
                  {{ d.status }}
                </span>
              </td>
              <td>
                <span v-if="d.distribution_round_id">
                  Tournée #{{ d.distribution_round_id }} (Bénévole: {{ d.distribution_round?.volunteer?.user?.firstname }} {{ d.distribution_round?.volunteer?.user?.last_name }})
                </span>
                <span v-else style="color: #ef4444; font-weight: 600;">
                Non assignée
                </span>
              </td>
              <td>
                <div style="display: flex; gap: 4px;">
                  <select v-model="selectedRoundForDelivery[d.ID]" style="font-size:0.8rem; padding: 4px;">
                    <option :value="null">-- Choisir une tournée --</option>
                    <option v-for="r in rounds" :key="r.ID" :value="r.ID">
                      Tournée #{{ r.ID }} ({{ r.volunteer?.user?.firstname }} {{ r.volunteer?.user?.last_name }} — {{ r.volunteer?.zone_area || 'Toutes zones' }})
                    </option>
                  </select>
                  <button class="btn btn-primary btn-sm" @click="assignDeliveryToRound(d.ID)">Assigner</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/api'
export default {
  name: 'StaffDashboard',
  data() {
    return {
      tab: 'merchants',
      msg: '', msgType: '',
      prefix: '/staff',
      pendingMerchants: [],
      volunteers: [],
      allSkills: [],
      skillAssign: {},
      newVolunteer: { firstname: '', lastname: '', email: '', password: '', zone_area: '', availability: '', vehicle: false },
      newSkill: { name: '', category: '' },
      newMission: { merchant_id: null, volunteer_id: null, product_id: null, pickup_date: '' },
      newService: { title: '', description: '', category: '', max_participants: 10, date: '', location: '', volunteer_id: null },
      servicesList: [],
      rounds: [],
      roundStatus: {},
      newRound: { date: '', volunteer_id: null, notes: '' },
      newDelivery: {},
      deliveries: [],
      selectedRoundForDelivery: {}
    }
  },
  mounted() { this.loadAll() },
  methods: {
    flash(text, type) { this.msg = text; this.msgType = type; setTimeout(() => this.msg = '', 4000) },
    formatDate(d) { return d ? new Date(d).toLocaleDateString('fr-FR') : '' },
    loadAll() { this.loadMerchants(); this.loadVolunteers(); this.loadSkills(); this.loadServices(); this.loadRounds(); this.loadDeliveries() },
    async loadDeliveries() {
      try {
        const r = await api.get(this.prefix + '/deliveries')
        this.deliveries = r.data.data || []
      } catch(e) { this.deliveries = [] }
    },
    async assignDeliveryToRound(deliveryId) {
      const roundId = this.selectedRoundForDelivery[deliveryId]
      if (!roundId) {
        this.flash('Veuillez choisir une tournée', 'err')
        return
      }
      try {
        await api.put(this.prefix + '/deliveries/' + deliveryId + '/assign', { distribution_round_id: roundId })
        this.flash('Livraison assignée à la tournée avec succès', 'ok')
        this.loadDeliveries()
        this.loadRounds()
      } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async loadMerchants() {
      try { const r = await api.get(this.prefix + '/merchants/pending'); this.pendingMerchants = r.data.data || [] } catch(e) { this.pendingMerchants = [] }
    },
    async approveMerchant(id) {
      try { await api.put(this.prefix + '/merchants/' + id + '/approve'); this.flash('Commerçant approuvé', 'ok'); this.loadMerchants() } catch(e) { this.flash('Erreur', 'err') }
    },
    async loadVolunteers() {
      try { const r = await api.get(this.prefix + '/volunteers'); this.volunteers = r.data.data || [] } catch(e) { this.volunteers = [] }
    },
    async registerVolunteer() {
      try { await api.post(this.prefix + '/volunteers', this.newVolunteer); this.flash('Bénévole inscrit', 'ok'); this.newVolunteer = { firstname: '', lastname: '', email: '', password: '', zone_area: '', availability: '', vehicle: false }; this.loadVolunteers() } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async loadSkills() {
      try { const r = await api.get(this.prefix + '/skills'); this.allSkills = r.data.data || [] } catch(e) { this.allSkills = [] }
    },
    async createSkill() {
      try { await api.post(this.prefix + '/skills', this.newSkill); this.flash('Compétence ajoutée', 'ok'); this.newSkill = { name: '', category: '' }; this.loadSkills() } catch(e) { this.flash('Erreur', 'err') }
    },
    async deleteSkill(id) {
      try { await api.delete(this.prefix + '/skills/' + id); this.flash('Supprimée', 'ok'); this.loadSkills() } catch(e) { this.flash('Erreur', 'err') }
    },
    async assignSkills(volunteerId) {
      const ids = this.skillAssign[volunteerId] || []
      if (!ids.length) return
      try { await api.put(this.prefix + '/volunteers/' + volunteerId + '/skills', { skill_ids: ids }); this.flash('Compétences assignées', 'ok'); this.loadVolunteers() } catch(e) { this.flash('Erreur', 'err') }
    },
    async createMission() {
      try {
        const payload = { ...this.newMission, pickup_date: new Date(this.newMission.pickup_date).toISOString() }
        await api.post(this.prefix + '/missions', payload)
        this.flash('Mission créée', 'ok')
        this.newMission = { merchant_id: null, volunteer_id: null, product_id: null, pickup_date: '' }
      } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async loadServices() {
      try { const r = await api.get(this.prefix + '/services'); this.servicesList = r.data.data || [] } catch(e) { this.servicesList = [] }
    },
    async createService() {
      try {
        const payload = { ...this.newService, date: new Date(this.newService.date).toISOString() }
        await api.post(this.prefix + '/services', payload)
        this.flash('Service créé', 'ok')
        this.newService = { title: '', description: '', category: '', max_participants: 10, date: '', location: '', volunteer_id: null }
        this.loadServices()
      } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async deleteService(id) {
      try { await api.delete(this.prefix + '/services/' + id); this.flash('Supprimé', 'ok'); this.loadServices() } catch(e) { this.flash('Erreur', 'err') }
    },
    async loadRounds() {
      try {
        const r = await api.get(this.prefix + '/distribution-rounds')
        this.rounds = r.data.data || []
        this.rounds.forEach(rd => { this.roundStatus[rd.ID] = rd.status })
      } catch(e) { this.rounds = [] }
    },
    async createRound() {
      try {
        const payload = { ...this.newRound, date: new Date(this.newRound.date).toISOString() }
        await api.post(this.prefix + '/distribution-rounds', payload)
        this.flash('Tournée créée', 'ok')
        this.newRound = { date: '', volunteer_id: null, notes: '' }
        this.loadRounds()
      } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async updateRoundStatus(id) {
      try { await api.put(this.prefix + '/distribution-rounds/' + id, { status: this.roundStatus[id] }); this.flash('Statut mis à jour', 'ok'); this.loadRounds() } catch(e) { this.flash('Erreur', 'err') }
    },
    async deleteRound(id) {
      try { await api.delete(this.prefix + '/distribution-rounds/' + id); this.flash('Tournée supprimée', 'ok'); this.loadRounds() } catch(e) { this.flash('Erreur', 'err') }
    },
    async addDelivery(roundId) {
      const d = this.newDelivery
      const payload = {
        recipient_name: d[roundId + '_name'],
        recipient_address: d[roundId + '_addr'],
        recipient_type: d[roundId + '_type'] || 'association',
        product_id: d[roundId + '_pid'],
        quantity: d[roundId + '_qty']
      }
      try { await api.post(this.prefix + '/distribution-rounds/' + roundId + '/deliveries', payload); this.flash('Livraison ajoutée', 'ok'); this.loadRounds() } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    downloadPDF(id) {
      const token = localStorage.getItem('token')
      window.open('http://localhost:8080' + this.prefix + '/distribution-rounds/' + id + '/pdf?token=' + token, '_blank')
    },
    exportServicesCSV() {
      const token = localStorage.getItem('token')
      window.open('http://localhost:8080' + this.prefix + '/services/export?token=' + token, '_blank')
    }
  }
}
</script>