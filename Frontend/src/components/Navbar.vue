<template>
  <nav class="navbar">
    <div class="navbar-brand">
      <svg class="brand-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M11 20A7 7 0 0 1 9.8 6.1C15.5 5 17 4.48 19 2c1 2 2 4.18 2 8 0 5.5-4.78 10-10 10Z"></path>
        <path d="M2 21c0-3 1.85-5.36 5.08-6C9.5 14.52 12 13 13 12"></path>
      </svg>
      <span class="brand-name">{{ $t('navbar.brand') }}</span>
    </div>
    <div class="navbar-info">
      
      <select v-model="currentLang" @change="changeLanguage" class="lang-select">
        <option value="fr">FR</option>
        <option value="en">EN</option>
        <option value="it">IT</option>
      </select>

      <template v-if="isLoggedIn">
        <span class="navbar-role-badge">{{ currentRole }}</span>
        <button class="btn btn-danger btn-sm" @click="logout">{{ $t('navbar.logout') }}</button>
      </template>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'Navbar',
  data() {
    return {
      currentRole: localStorage.getItem('role') || '',
      isLoggedIn: !!localStorage.getItem('token'),
      currentLang: localStorage.getItem('lang') || 'fr'
    }
  },
  watch: {
    '$route'() {
      this.currentRole = localStorage.getItem('role') || ''
      this.isLoggedIn = !!localStorage.getItem('token')
    }
  },
  methods: {
    changeLanguage() {
      this.$i18n.locale = this.currentLang
      localStorage.setItem('lang', this.currentLang)
    },
    logout() {
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      this.isLoggedIn = false
      this.currentRole = ''
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2.5rem;
  background: var(--bg-dark, #0f172a);
  color: #ffffff;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  font-family: 'Outfit', 'Inter', sans-serif;
  font-size: 1.3rem;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.brand-icon {
  width: 24px;
  height: 24px;
  color: var(--primary-color, #2d6a4f);
}

.navbar-info {
  display: flex;
  align-items: center;
  gap: 1.2rem;
}

.navbar-role-badge {
  background: rgba(45, 106, 79, 0.25);
  color: #52b788;
  border: 1px solid rgba(82, 183, 136, 0.4);
  padding: 0.35rem 0.9rem;
  border-radius: 9999px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.lang-select {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 0.35rem 0.7rem;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  outline: none;
  transition: all 0.2s ease;
}

.lang-select:hover {
  background: rgba(255, 255, 255, 0.2);
}

.lang-select option {
  background: #1e293b;
  color: #fff;
}
</style>