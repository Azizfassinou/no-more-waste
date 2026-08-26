<template>
  <div class="dashboard">
    <h1>{{ $t('merchant.title') }}</h1>

    <div class="tabs">
      <button :class="{ active: tab === 'my-products' }" @click="tab = 'my-products'">{{ $t('merchant.tabs.myProducts') }}</button>
      <button :class="{ active: tab === 'add-product' }" @click="tab = 'add-product'">{{ $t('merchant.tabs.addProduct') }}</button>
      <button :class="{ active: tab === 'profile' }" @click="tab = 'profile'">{{ $t('merchant.tabs.profile') }}</button>
    </div>

    <div v-if="msg" :class="'alert ' + (msgType === 'ok' ? 'alert-success' : 'alert-error')">{{ msg }}</div>

    <div v-if="renewalDate && (isSubscriptionExpired || isSubscriptionExpiringSoon)" class="subscription-banner" :class="isSubscriptionExpired ? 'banner-expired' : 'banner-warning'">
        <strong v-if="isSubscriptionExpired">Votre abonnement annuel a expiré le {{ formattedRenewalDate }}.</strong>
        <strong v-else>Votre abonnement expire dans {{ daysUntilRenewal }} jour(s) (le {{ formattedRenewalDate }}).</strong>
        <p style="margin: 0.2rem 0 0 0; font-size: 0.88rem; opacity: 0.9;">
          Renouvelez dès maintenant pour continuer à publier et vendre vos invendus sans interruption.
        </p>
        <button class="btn-renew-banner" @click="renewSubscription">
        Renouveler (15,00 €)
        </button>
    </div>

    <div v-if="tab === 'my-products'">
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1rem; margin-bottom: 2rem;">
        <div style="background: linear-gradient(135deg, #1e3a8a, #3b82f6); color: white; padding: 1.2rem; border-radius: 12px; box-shadow: 0 4px 12px rgba(30, 58, 138, 0.2);">
          <div style="font-size: 0.82rem; text-transform: uppercase; letter-spacing: 0.5px; opacity: 0.9;">Ma Cagnotte (Chiffre d'Affaires)</div>
          <div style="font-size: 1.8rem; font-weight: 700; margin-top: 0.4rem;">{{ stats.cagnotte_revenue ? stats.cagnotte_revenue.toFixed(2) : '0.00' }} €</div>
        </div>

        <div style="background: linear-gradient(135deg, #0f766e, #0d9488); color: white; padding: 1.2rem; border-radius: 12px; box-shadow: 0 4px 12px rgba(13, 148, 136, 0.2);">
          <div style="font-size: 0.82rem; text-transform: uppercase; letter-spacing: 0.5px; opacity: 0.9;">Fin d'Abonnement</div>
          <div style="font-size: 1.25rem; font-weight: 700; margin-top: 0.4rem; display: flex; align-items: center; justify-content: space-between;">
            <span>{{ formattedRenewalDate }}</span>
            <span v-if="isSubscriptionExpired" class="badge-status badge-status-danger">Expiré</span>
            <span v-else-if="isSubscriptionExpiringSoon" class="badge-status badge-status-warning">Expire bientôt</span>
            <span v-else class="badge-status badge-status-success">Actif</span>
          </div>
          <button style="margin-top: 0.6rem; background: rgba(255,255,255,0.2); border: none; color: white; padding: 0.35rem 0.75rem; border-radius: 6px; font-size: 0.78rem; font-weight: 600; cursor: pointer; transition: background 0.2s;" @click="renewSubscription">
            🔄 Renouveler (15 €)
          </button>
        </div>

        <div style="background: white; padding: 1.2rem; border-radius: 12px; border: 1px solid #e2e8f0; box-shadow: 0 1px 3px rgba(0,0,0,0.05);">
          <div style="font-size: 0.88rem; color: #64748b;">Produits Publiés</div>
          <div style="font-size: 1.8rem; font-weight: 700; color: #1e293b; margin-top: 0.4rem;">{{ stats.total_published_products || 0 }}</div>
        </div>

        <div style="background: white; padding: 1.2rem; border-radius: 12px; border: 1px solid #e2e8f0; box-shadow: 0 1px 3px rgba(0,0,0,0.05);">
          <div style="font-size: 0.88rem; color: #64748b;">Stock Restant</div>
          <div style="font-size: 1.8rem; font-weight: 700; color: #d97706; margin-top: 0.4rem;">{{ stats.total_remaining_stock || 0 }}</div>
        </div>
      </div>

      <h2 class="section-title">{{ $t('merchant.myProductsTitle') }}</h2>
      <div v-if="products.length === 0" class="empty-state">{{ $t('merchant.noProducts') }}</div>
      <div class="card-grid">
        <div class="card" v-for="p in products" :key="p.id || p.ID">
          <h3>{{ p.title }}</h3>
          <p>{{ p.description }}</p>
          <p><strong>{{ $t('client.barcode') }} :</strong> <code>{{ p.barcode }}</code></p>
          <p><strong>Prix d'origine :</strong> {{ p.original_price }}€ → <strong>Prix réduis :</strong> <span style="color:#2d6a4f; font-weight:600;">{{ p.discount_price }}€</span></p>
          <p><strong>Stock Restant :</strong> <span style="font-weight:700; color:#d97706;">{{ p.quantity }}</span> | <strong>Unités Vendues :</strong> <span style="font-weight:700; color:#2d6a4f;">{{ p.sold_units || 0 }}</span></p>
          <p><strong>Revenu Généré :</strong> <span style="font-weight:700; color:#1e3a8a;">{{ p.revenue ? p.revenue.toFixed(2) : '0.00' }} €</span></p>
          <p><strong>{{ $t('client.expiresOn') }} :</strong> {{ formatDate(p.expiry_date) }}</p>
          <span class="badge" :class="p.is_available ? 'badge-green' : 'badge-red'">{{ p.is_available ? $t('merchant.available') : $t('merchant.soldOut') }}</span>
        </div>
      </div>
    </div>

    <div v-if="tab === 'add-product'">
      <h2 class="section-title">{{ $t('merchant.addProductTitle') }}</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('merchant.productTitle') }}</label>
            <input v-model="newProduct.title" placeholder="Panier fruits & légumes" />
          </div>
          <div class="form-group">
            <label>{{ $t('merchant.productDescription') }}</label>
            <input v-model="newProduct.description" placeholder="Lot de fruits invendus" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('client.originalPrice') }} (€)</label>
            <input v-model.number="newProduct.original_price" type="number" step="0.01" min="0" />
          </div>
          <div class="form-group">
            <label>{{ $t('client.discountPrice') }} (€)</label>
            <input v-model.number="newProduct.discount_price" type="number" step="0.01" min="0" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('client.quantity') }}</label>
            <input v-model.number="newProduct.quantity" type="number" min="1" />
          </div>
          <div class="form-group">
            <label>{{ $t('merchant.expiryDate') }}</label>
            <input v-model="newProduct.expiry_date" type="date" />
          </div>
        </div>
        <button class="btn btn-primary" @click="createProduct">{{ $t('merchant.publish') }}</button>
        <p style="margin-top:0.5rem;font-size:0.8rem;color:#888;">{{ $t('merchant.autoBarcodeHint') }}</p>
      </div>
    </div>

    <div v-if="tab === 'profile'">
      <h2 class="section-title">{{ $t('merchant.profileTitle') }}</h2>
      <div class="form-section">
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('auth.firstname') }}</label>
            <input v-model="merchantProfile.firstname" />
          </div>
          <div class="form-group">
            <label>{{ $t('auth.lastname') }}</label>
            <input v-model="merchantProfile.last_name" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('auth.phone') }}</label>
            <input v-model="merchantProfile.phone" />
          </div>
          <div class="form-group">
            <label>{{ $t('merchant.personalAddress') }}</label>
            <input v-model="merchantProfile.address" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>{{ $t('auth.companyName') }}</label>
            <input v-model="merchantProfile.company_name" />
          </div>
          <div class="form-group">
            <label>{{ $t('merchant.companyAddress') }}</label>
            <input v-model="merchantProfile.company_address" />
          </div>
        </div>
        <button class="btn btn-primary" @click="updateMerchantProfile">{{ $t('client.updateProfile') }}</button>
      </div>

      <h2 class="section-title" style="margin-top: 2rem;">Statut de l'Abonnement Commerçant</h2>
      <div class="form-section subscription-card">
        <div class="subscription-details">
          <div class="subscription-item">
            <span class="sub-label">Date d'échéance / Renouvellement :</span>
            <span class="sub-value font-highlight">{{ formattedRenewalDate }}</span>
          </div>
          <div class="subscription-item">
            <span class="sub-label">Statut actuel :</span>
            <span v-if="isSubscriptionExpired" class="badge-status badge-status-danger">Expiré</span>
            <span v-else-if="isSubscriptionExpiringSoon" class="badge-status badge-status-warning">Expire sous {{ daysUntilRenewal }} jours</span>
            <span v-else class="badge-status badge-status-success">Abonnement Actif (Valide)</span>
          </div>
          <div class="subscription-item">
            <span class="sub-label">Tarif annuel :</span>
            <span class="sub-value">15,00 € / an (1ère année offerte)</span>
          </div>
        </div>
        <p style="font-size: 0.9rem; color: #64748b; margin-top: 1rem; margin-bottom: 1.2rem; line-height: 1.5;">
          Le renouvellement prolonge votre accès d'un an supplémentaire à compter de la date d'échéance. Le paiement est sécurisé par Stripe.
        </p>
        <button class="btn btn-primary btn-renew-stripe" @click="renewSubscription">
          Renouveler mon abonnement annuel (15,00 €)
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/api'

export default {
  name: 'MerchantDashboard',
  data() {
    return {
      tab: 'my-products',
      msg: '', msgType: '',
      products: [],
      stats: { cagnotte_revenue: 0, total_published_products: 0, total_sold_units: 0, total_remaining_stock: 0, renewal_date: null },
      newProduct: { title: '', description: '', original_price: 0, discount_price: 0, quantity: 1, expiry_date: '' },
      merchantProfile: {}
    }
  },
  computed: {
    renewalDate() {
      return this.stats.renewal_date || this.merchantProfile.renewal_date || null
    },
    formattedRenewalDate() {
      if (!this.renewalDate) return 'Non renseignée'
      return new Date(this.renewalDate).toLocaleDateString('fr-FR', {
        day: '2-digit',
        month: 'long',
        year: 'numeric'
      })
    },
    daysUntilRenewal() {
      if (!this.renewalDate) return 999
      const now = new Date()
      const exp = new Date(this.renewalDate)
      return Math.ceil((exp - now) / (1000 * 60 * 60 * 24))
    },
    isSubscriptionExpired() {
      return this.daysUntilRenewal < 0
    },
    isSubscriptionExpiringSoon() {
      return this.daysUntilRenewal >= 0 && this.daysUntilRenewal <= 30
    }
  },
  mounted() { 
    this.loadProducts(); 
    this.loadProfile();
    
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('payment') === 'success') {
      this.flash('Paiement Stripe validé ! Votre abonnement commerçant est renouvelé avec succès.', 'ok');
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  },
  methods: {
    flash(text, type) { this.msg = text; this.msgType = type; setTimeout(() => this.msg = '', 5000) },
    formatDate(d) { return d ? new Date(d).toLocaleDateString('fr-FR') : '' },
    async loadProducts() {
      try {
        const r = await api.get('/merchant/stats')
        this.stats = r.data || {}
        this.products = r.data.products || []
        if (r.data.renewal_date) {
          this.merchantProfile.renewal_date = r.data.renewal_date
        }
      } catch(e) {
        this.products = []
      }
    },
    async createProduct() {
      if (this.newProduct.original_price < 0 || this.newProduct.discount_price < 0) {
        this.flash("Le prix ne peut pas être inférieur à 0 €.", 'err')
        return
      }
      if (this.newProduct.quantity <= 0) {
        this.flash("La quantité doit être supérieure à 0.", 'err')
        return
      }
      if (this.newProduct.discount_price > this.newProduct.original_price) {
        this.flash("Le prix réduit ne peut pas être supérieur au prix d'origine.", 'err')
        return
      }
      try {
        const payload = { ...this.newProduct, expiry_date: new Date(this.newProduct.expiry_date).toISOString() }
        const r = await api.post('/merchant/product', payload)
        this.flash('Produit publié ! Code-barre : ' + r.data.data.barcode, 'ok')
        this.newProduct = { title: '', description: '', original_price: 0, discount_price: 0, quantity: 1, expiry_date: '' }
        this.loadProducts()
      } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async loadProfile() {
      try {
        const r = await api.get('/profile');
        this.merchantProfile = r.data.data || {};
        if (r.data.merchant) {
          this.merchantProfile.renewal_date = r.data.merchant.RenewalDate || r.data.merchant.renewal_date;
          this.merchantProfile.company_name = r.data.merchant.CompanyName || r.data.merchant.company_name;
          this.merchantProfile.company_address = r.data.merchant.CompanyAddress || r.data.merchant.company_address;
          this.merchantProfile.merchant_id = r.data.merchant.ID || r.data.merchant.id;
        }
      } catch(e) {}
    },
    async updateMerchantProfile() {
      try { await api.put('/merchant/' + (this.merchantProfile.merchant_id || this.merchantProfile.ID), this.merchantProfile); this.flash('Profil mis à jour', 'ok') } catch(e) { this.flash(e.response?.data?.error || 'Erreur', 'err') }
    },
    async renewSubscription() {
      try {
        const r = await api.post('/merchant/subscription/renew');
        if (r.data.url) {
          window.location.href = r.data.url;
        } else {
          this.flash(r.data.message || 'Abonnement renouvelé avec succès.', 'ok');
          if (r.data.renewal_date) {
            this.stats.renewal_date = r.data.renewal_date;
            this.merchantProfile.renewal_date = r.data.renewal_date;
          }
        }
      } catch(e) {
        this.flash(e.response?.data?.error || 'Erreur lors du renouvellement', 'err');
      }
    }
  }
}
</script>

<style scoped>
.subscription-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.1rem 1.5rem;
  border-radius: 12px;
  margin-bottom: 1.5rem;
  gap: 1rem;
  flex-wrap: wrap;
}

.banner-warning {
  background-color: #fffbebfb;
  border: 1px solid #fcd34d;
  color: #92400e;
}

.banner-expired {
  background-color: #fef2f2;
  border: 1px solid #fca5a5;
  color: #991b1b;
}

.banner-content {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
}

.banner-icon {
  font-size: 1.8rem;
}

.btn-renew-banner {
  background: #2d6a4f;
  color: white;
  border: none;
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 4px rgba(45, 106, 79, 0.2);
}

.btn-renew-banner:hover {
  background: #1b4332;
}

.subscription-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
}

.subscription-details {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.subscription-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.95rem;
}

.sub-label {
  font-weight: 600;
  color: #475569;
  min-width: 220px;
}

.sub-value {
  color: #1e293b;
  font-weight: 500;
}

.font-highlight {
  font-weight: 700;
  color: #0f766e;
  font-size: 1.05rem;
}

.badge-status {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.badge-status-success {
  background-color: #dcfce7;
  color: #166534;
}

.badge-status-warning {
  background-color: #fef3c7;
  color: #92400e;
}

.badge-status-danger {
  background-color: #fee2e2;
  color: #991b1b;
}

.btn-renew-stripe {
  background-color: #6366f1;
  border-color: #6366f1;
  font-size: 0.95rem;
  padding: 0.75rem 1.5rem;
}

.btn-renew-stripe:hover {
  background-color: #4f46e5;
}
</style>
