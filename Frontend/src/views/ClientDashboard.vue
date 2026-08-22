<template>
  <div class="dashboard">
    <h1>{{ $t('client.title') }}</h1>

    <div class="tabs">
      <button :class="{ active: tab === 'products' }" @click="tab = 'products'">{{ $t('client.tabs.products') }}</button>
      <button :class="{ active: tab === 'cart' }" @click="tab = 'cart'">
        Mon Panier ({{ cartCount }}) - {{ cartTotal.toFixed(2) }}€
      </button>
      <button :class="{ active: tab === 'search' }" @click="tab = 'search'">{{ $t('client.tabs.search') }}</button>
      <button :class="{ active: tab === 'services' }" @click="tab = 'services'">{{ $t('client.tabs.services') }}</button>
      <button :class="{ active: tab === 'my-services' }" @click="tab = 'my-services'">{{ $t('client.tabs.myServices') }}</button>
      <button :class="{ active: tab === 'orders' }" @click="tab = 'orders'">Mes Commandes</button>
      <button :class="{ active: tab === 'profile' }" @click="tab = 'profile'">{{ $t('client.tabs.profile') }}</button>
    </div>

    <div v-if="msg" :class="'alert ' + (msgType === 'ok' ? 'alert-success' : 'alert-error')">{{ msg }}</div>

    <div v-if="tab === 'products'">
      <h2 class="section-title">{{ $t('client.availableProducts') }}</h2>
      <div v-if="products.length === 0" class="empty-state">{{ $t('client.noProducts') }}</div>
      <div class="card-grid">
        <div class="card" v-for="p in products" :key="p.ID">
          <h3>{{ p.title }}</h3>
          <p>{{ p.description }}</p>
          <p><strong>{{ $t('client.barcode') }} :</strong> {{ p.barcode }}</p>
          <p><strong>{{ $t('client.price') }} :</strong> <s>{{ p.original_price }}€</s> → <span style="color:#2d6a4f;font-weight:600;">{{ p.discount_price }}€</span></p>
          <p><strong>{{ $t('client.quantity') }} :</strong> {{ p.quantity }} | <strong>{{ $t('client.expiresOn') }} :</strong> {{ formatDate(p.expiry_date) }}</p>
          <div class="btn-group">
            <button class="btn btn-primary btn-sm" @click="addToCart(p)">Ajouter au Panier</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'cart'">
      <h2 class="section-title">Mon Panier d'Achat</h2>
      <div v-if="cart.length === 0" class="empty-state">Votre panier est vide. Parcourez nos produits et ajoutez des articles !</div>
      <div v-else>
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Produit</th>
                <th>Prix Unitaire</th>
                <th>Quantité</th>
                <th>Sous-total</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in cart" :key="item.ID">
                <td><strong>{{ item.title }}</strong></td>
                <td>{{ item.discount_price.toFixed(2) }} €</td>
                <td>
                  <input type="number" min="1" :max="item.max_quantity" v-model.number="item.quantity" style="width:60px; padding:4px;" />
                </td>
                <td><strong>{{ (item.discount_price * item.quantity).toFixed(2) }} €</strong></td>
                <td>
                  <button class="btn btn-danger btn-sm" @click="removeFromCart(item.ID)">Retirer</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div style="margin-top: 1.5rem; background: #fff; padding: 1.2rem; border-radius: 8px; border: 1px solid #cbd5e1;">
          <h3 style="margin-bottom: 0.5rem;">Récapitulatif de la commande</h3>
          <p style="font-size: 1.2rem; font-weight: 700; color: #1e293b; margin-bottom: 1rem;">
            Total du panier : <span style="color: #2d6a4f;">{{ cartTotal.toFixed(2) }} €</span>
          </p>

          <div v-if="cartTotal < 10.0" class="alert alert-error" style="margin-bottom: 1rem; font-weight: 600;">
            Montant minimum requis : 10,00 € pour passer commande et procéder au paiement Stripe. 
            (Il vous manque {{ (10.0 - cartTotal).toFixed(2) }} €)
          </div>
          <div v-else class="alert alert-success" style="margin-bottom: 1rem; font-weight: 600;">
            Montant minimum atteint ({{ cartTotal.toFixed(2) }} €). Vous pouvez procéder au paiement Stripe !
          </div>

          <div class="form-group" style="margin-bottom: 1.2rem;">
            <label style="font-weight:600;">Adresse de livraison :</label>
            <input v-model="deliveryAddress" placeholder="12 Rue de la Paix, 75002 Paris" style="padding: 0.5rem; border-radius:6px; border:1px solid #cbd5e1; width:100%;" />
          </div>

          <button 
            class="btn btn-primary btn-block" 
            :disabled="cartTotal < 10.0 || cart.length === 0 || loading" 
            @click="checkoutCart"
            style="font-size: 1.05rem; padding: 0.8rem;"
          >
             {{ loading ? 'Création de la session Stripe...' : 'Payer par Stripe (' + cartTotal.toFixed(2) + ' €)' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="tab === 'search'">
      <h2 class="section-title">{{ $t('client.searchTitle') }}</h2>
      <div class="form-section">
        <div class="form-group">
          <label>{{ $t('client.searchKeywords') }}</label>
          <input v-model="searchQuery" :placeholder="$t('client.searchPlaceholder')" @keyup.enter="searchProducts" />
        </div>
        <div class="btn-group">
          <button class="btn btn-primary" @click="searchProducts">{{ $t('common.search') }}</button>
        </div>
      </div>
      <div v-if="searchResults.length" class="card-grid" style="margin-top:1rem;">
        <div class="card" v-for="p in searchResults" :key="p.ID">
          <h3>{{ p.title }}</h3>
          <p>{{ p.description }}</p>
          <p><strong>{{ $t('client.discountPrice') }} :</strong> {{ p.discount_price }}€ | <strong>{{ $t('client.quantity') }} :</strong> {{ p.quantity }}</p>
        </div>
      </div>
    </div>

    <div v-if="tab === 'services'">
      <h2 class="section-title">{{ $t('client.availableServicesTitle') }}</h2>
      <div v-if="servicesList.length === 0" class="empty-state">{{ $t('client.noServices') }}</div>
      <div class="card-grid">
        <div class="card" v-for="s in servicesList" :key="s.ID">
          <h3>{{ s.title }}</h3>
          <p>{{ s.description }}</p>
          <p><span class="badge badge-blue">{{ s.category }}</span> <span class="badge" :class="s.status === 'open' ? 'badge-green' : 'badge-orange'">{{ s.status }}</span></p>
          <p><strong>{{ $t('common.date') }} :</strong> {{ formatDate(s.date) }} | <strong>{{ $t('common.location') }} :</strong> {{ s.location }}</p>
          <p><strong>{{ $t('client.maxPlaces') }} :</strong> {{ s.max_participants }}</p>
          <div class="btn-group">
            <button class="btn btn-primary btn-sm" @click="registerToService(s.ID)" v-if="s.status === 'open'">{{ $t('client.register') }}</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'my-services'">
      <h2 class="section-title">{{ $t('client.myRegistrationsTitle') }}</h2>
      <div v-if="myRegistrations.length === 0" class="empty-state">{{ $t('client.noRegistrations') }}</div>
      <div class="card-grid">
        <div class="card" v-for="r in myRegistrations" :key="r.ID">
          <h3>{{ r.service.title }}</h3>
          <p>{{ r.service.description }}</p>
          <p><strong>{{ $t('common.date') }} :</strong> {{ formatDate(r.service.date) }} | <strong>{{ $t('common.location') }} :</strong> {{ r.service.location }}</p>
          <div class="btn-group">
            <button class="btn btn-danger btn-sm" @click="unregisterFromService(r.service.ID)">{{ $t('client.unregister') }}</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'profile'">
      <h2 class="section-title">{{ $t('client.profileTitle') }}</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('auth.firstname') }}</label>
            <input v-model="profile.firstname" />
          </div>
          <div class="form-group">
            <label>{{ $t('auth.lastname') }}</label>
            <input v-model="profile.last_name" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('auth.email') }}</label>
            <input v-model="profile.email" type="email" />
          </div>
          <div class="form-group">
            <label>{{ $t('auth.phone') }}</label>
            <input v-model="profile.phone" />
          </div>
        </div>
        <div class="form-group">
          <label>{{ $t('auth.address') }}</label>
          <input v-model="profile.address" />
        </div>
        <button class="btn btn-primary" @click="updateProfile">{{ $t('client.updateProfile') }}</button>
      </div>
    </div>

    <div v-if="tab === 'orders'">
      <h2 class="section-title">Historique de Mes Commandes Validées</h2>
      <div v-if="myOrders.length === 0" class="empty-state">Vous n'avez pas encore passé de commande.</div>
      <div class="table-container" v-else>
        <table>
          <thead>
            <tr>
              <th>ID Commande</th>
              <th>Produit</th>
              <th>Commerçant</th>
              <th>Montant Total</th>
              <th>Date de Commande</th>
              <th>Statut Paiement</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="o in myOrders" :key="o.ID">
              <td>#{{ o.ID }}</td>
              <td><strong>{{ o.product?.title || 'Produit #' + o.product_id }}</strong></td>
              <td>{{ o.product?.merchant?.company_name || 'Partenaire No-More-Waste' }}</td>
              <td><strong>{{ o.total_price.toFixed(2) }} €</strong></td>
              <td>{{ formatDate(o.CreatedAt) }}</td>
              <td>
                <span class="badge badge-green">Paiement Validé (Stripe)</span>
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
  name: 'ClientDashboard',
  data() {
    return {
      tab: 'products',
      msg: '', msgType: '',
      loading: false,
      products: [],
      cart: [],
      deliveryAddress: '',
      searchQuery: '', barcodeQuery: '',
      searchResults: [], scannedProduct: null,
      servicesList: [],
      myRegistrations: [],
      myOrders: [],
      profile: {}
    }
  },
  computed: {
    cartTotal() {
      return this.cart.reduce((sum, item) => sum + (item.discount_price * item.quantity), 0)
    },
    cartCount() {
      return this.cart.reduce((sum, item) => sum + item.quantity, 0)
    }
  },
  mounted() { 
    this.loadProducts(); 
    this.loadServices(); 
    this.loadMyRegistrations(); 
    this.loadMyOrders();
    this.loadProfile();

    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('payment') === 'success') {
      this.flash('Paiement Stripe validé ! Commande confirmée.', 'ok');
      this.cart = [];
      this.loadMyOrders();
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  },
  methods: {
    flash(text, type) { this.msg = text; this.msgType = type; setTimeout(() => this.msg = '', 5000) },
    formatDate(d) { return d ? new Date(d).toLocaleDateString('fr-FR') : '' },
    async loadProducts() {
      try { const r = await api.get('/client/products'); this.products = r.data.data || [] } catch(e) { this.products = [] }
    },
    async loadMyOrders() {
      try { const r = await api.get('/client/my-orders'); this.myOrders = r.data.data || [] } catch(e) { this.myOrders = [] }
    },
    addToCart(product) {
      const existing = this.cart.find(item => item.ID === product.ID)
      if (existing) {
        if (existing.quantity < product.quantity) {
          existing.quantity++
          this.flash('Quantité mise à jour dans votre panier', 'ok')
        } else {
          this.flash('Stock maximum atteint pour ce produit', 'err')
        }
      } else {
        this.cart.push({
          ID: product.ID,
          title: product.title,
          discount_price: product.discount_price,
          quantity: 1,
          max_quantity: product.quantity
        })
        this.flash('Produit ajouté à votre panier !', 'ok')
      }
    },
    removeFromCart(id) {
      this.cart = this.cart.filter(item => item.ID !== id)
      this.flash('Produit retiré du panier', 'ok')
    },
    async checkoutCart() {
      if (this.cartTotal < 10.0) {
        this.flash('Le montant minimum de la commande est de 10,00 € pour procéder au paiement Stripe.', 'err')
        return
      }
      this.loading = true
      try {
        const payload = {
          items: this.cart.map(i => ({ product_id: i.ID, quantity: i.quantity })),
          delivery_address: this.deliveryAddress || this.profile.address
        }
        const res = await api.post('/client/orders', payload)
        if (res.data.url) {
          window.location.href = res.data.url;
        } else {
          this.flash(res.data.message || 'Paiement Stripe réussi et commande validée !', 'ok')
          this.cart = []
          this.tab = 'products'
          this.loadProducts()
        }
      } catch (e) {
        this.flash(e.response?.data?.error || 'Erreur lors de la validation', 'err')
      } finally {
        this.loading = false
      }
    },
    async searchProducts() {
      if (!this.searchQuery) return
      try { const r = await api.get('/client/products/search?q=' + this.searchQuery); this.searchResults = r.data.data || [] } catch(e) { this.flash('Erreur de recherche', 'err') }
    },
    async scanBarcode() {
      if (!this.barcodeQuery) return
      try { const r = await api.get('/client/products/scan?code=' + this.barcodeQuery); this.scannedProduct = r.data.data } catch(e) { this.flash('Produit non trouvé', 'err'); this.scannedProduct = null }
    },
    async loadServices() {
      try { const r = await api.get('/client/services'); this.servicesList = r.data.data || [] } catch(e) { this.servicesList = [] }
    },
    async registerToService(id) {
      try { await api.post('/client/services/' + id + '/register'); this.flash('Inscrit !', 'ok'); this.loadServices(); this.loadMyRegistrations() } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async loadMyRegistrations() {
      try { const r = await api.get('/client/my-services'); this.myRegistrations = r.data.data || [] } catch(e) { this.myRegistrations = [] }
    },
    async unregisterFromService(id) {
      try { await api.delete('/client/services/' + id + '/register'); this.flash('Désinscrit', 'ok'); this.loadMyRegistrations(); this.loadServices() } catch(e) { this.flash('Erreur', 'err') }
    },
    async loadProfile() {
      try {
        const r = await api.get('/profile')
        this.profile = r.data.data || {}
        if (!this.deliveryAddress && this.profile.address) {
          this.deliveryAddress = this.profile.address
        }
      } catch(e) {}
    },
    async updateProfile() {
      try { await api.put('/profile', this.profile); this.flash('Profil mis à jour', 'ok') } catch(e) { this.flash('Erreur', 'err') }
    }
  }
}
</script>