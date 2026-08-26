<template>
  <div class="landing-page">
    <main class="hero-section">
      <div class="hero-container">
        
        <div class="hero-content">
          <h1 class="hero-title">
            {{ $t('landing.heroTitle') }} <br/>
            <span class="gradient-text">{{ $t('landing.heroSubtitle') }}</span>
          </h1>
          <p class="hero-description">
            {{ $t('landing.heroDescription') }}
          </p>

          <div class="pillars-grid">
            <div class="pillar-card">
              <div class="pillar-icon-box green">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#52b788" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
              </div>
              <div>
                <h3 class="pillar-title">{{ $t('landing.pillar1Title') }}</h3>
                <p class="pillar-desc">{{ $t('landing.pillar1Desc') }}</p>
              </div>
            </div>

            <div class="pillar-card">
              <div class="pillar-icon-box orange">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="21" r="1"></circle><circle cx="20" cy="21" r="1"></circle><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path></svg>
              </div>
              <div>
                <h3 class="pillar-title">{{ $t('landing.pillar2Title') }}</h3>
                <p class="pillar-desc">{{ $t('landing.pillar2Desc') }}</p>
              </div>
            </div>

            <div class="pillar-card">
              <div class="pillar-icon-box blue">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#60a5fa" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="1" y="3" width="15" height="13"></rect><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"></polygon><circle cx="5.5" cy="18.5" r="2.5"></circle><circle cx="18.5" cy="18.5" r="2.5"></circle></svg>
              </div>
              <div>
                <h3 class="pillar-title">{{ $t('landing.pillar3Title') }}</h3>
                <p class="pillar-desc">{{ $t('landing.pillar3Desc') }}</p>
              </div>
            </div>
          </div>

          <div class="impact-metrics" style="justify-content: center;">
            <div class="metric-item">
              <span class="metric-value">100%</span>
              <span class="metric-label">{{ $t('landing.metric1Label') }}</span>
            </div>
            <div class="metric-divider"></div>
            <div class="metric-item">
              <span class="metric-value">+50000</span>
              <span class="metric-label">{{ $t('landing.metric2Label') }}</span>
            </div>
          </div>
        </div>

        <div class="hero-form-wrapper">
          <div class="login-card">
            <div class="login-header">
              <h2>{{ $t('landing.portalTitle') }}</h2>
              <p class="login-subtitle">{{ $t('landing.portalSubtitle') }}</p>
              
              <div class="auth-tabs">
                <button :class="{ active: mode === 'login' }" @click="mode = 'login'">
                  {{ $t('auth.loginTab') }}
                </button>
                <button :class="{ active: mode === 'reg-client' }" @click="mode = 'reg-client'">
                  {{ $t('auth.registerClientTab') }}
                </button>
                <button :class="{ active: mode === 'reg-merchant' }" @click="mode = 'reg-merchant'">
                  {{ $t('auth.registerMerchantTab') }}
                </button>
                <button :class="{ active: mode === 'contact-volunteer' }" @click="mode = 'contact-volunteer'">
                  {{ $t('auth.volunteerTab') }}
                </button>
              </div>
            </div>

            <div v-if="msg" :class="'alert ' + (msgType === 'ok' ? 'alert-success' : 'alert-error')">{{ msg }}</div>

            <form v-if="mode === 'login'" @submit.prevent="handleLogin" class="auth-form">
              <div class="form-group">
                <label>{{ $t('auth.email') }}</label>
                <input type="email" v-model="email" required placeholder="exemple@domaine.com" />
              </div>
              <div class="form-group">
                <label>{{ $t('auth.password') }}</label>
                <div class="password-wrapper">
                  <input :type="showLoginPassword ? 'text' : 'password'" v-model="password" required placeholder="••••••••" />
                  <button type="button" class="toggle-password-btn" @click="showLoginPassword = !showLoginPassword" :title="showLoginPassword ? 'Masquer' : 'Afficher'">
                    <svg v-if="!showLoginPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                  </button>
                </div>
              </div>
              <div style="text-align: right; margin-top: -0.3rem; margin-bottom: 1.2rem;">
                <a href="#" @click.prevent="mode = 'forgot-password'; forgotStep = 1; forgotEmail = email;" style="font-size: 0.85rem; color: #2d6a4f; font-weight: 600; text-decoration: underline;">
                  {{ $t('auth.forgotPassword') }}
                </a>
              </div>
              <button type="submit" :disabled="loading" class="btn btn-primary btn-block">
                {{ loading ? $t('auth.loading') : $t('auth.submitLogin') }}
              </button>
            </form>

            <form v-if="mode === 'forgot-password' && forgotStep === 1" @submit.prevent="handleForgotPassword" class="auth-form">
              <h4 style="margin-bottom: 0.8rem; color: #1e293b; font-family: 'Outfit', sans-serif; font-size: 1.2rem;">{{ $t('auth.forgotPasswordTitle') }}</h4>
              <p style="font-size: 0.88rem; color: #64748b; margin-bottom: 1.2rem; line-height: 1.5;">
                {{ $t('auth.forgotPasswordInstruction') }}
              </p>
              <div class="form-group">
                <label>{{ $t('auth.accountEmail') }}</label>
                <input type="email" v-model="forgotEmail" required placeholder="votre.email@domaine.com" />
              </div>
              <button type="submit" :disabled="loading" class="btn btn-primary btn-block">
                {{ loading ? $t('auth.loading') : $t('auth.sendResetCode') }}
              </button>
              <div style="text-align: center; margin-top: 1rem;">
                <a href="#" @click.prevent="mode = 'login'" style="font-size: 0.85rem; color: #64748b; text-decoration: underline;">
                  {{ $t('auth.backToLogin') }}
                </a>
              </div>
            </form>

            <form v-if="mode === 'forgot-password' && forgotStep === 2" @submit.prevent="handleResetPassword" class="auth-form">
              <h4 style="margin-bottom: 0.8rem; color: #1e293b; font-family: 'Outfit', sans-serif; font-size: 1.2rem;">{{ $t('auth.resetPasswordTitle') }}</h4>
              <p style="font-size: 0.88rem; color: #64748b; margin-bottom: 1.2rem; line-height: 1.5;">
                {{ $t('auth.resetPasswordInstruction') }}
              </p>
              <div class="form-group">
                <label>{{ $t('auth.resetCodeLabel') }}</label>
                <input type="text" v-model="resetCode" required placeholder="123456" />
              </div>
              <div class="form-group">
                <label>{{ $t('auth.newPasswordLabel') }}</label>
                <div class="password-wrapper">
                  <input :type="showResetNewPassword ? 'text' : 'password'" v-model="newPassword" required placeholder="••••••••" />
                  <button type="button" class="toggle-password-btn" @click="showResetNewPassword = !showResetNewPassword" :title="showResetNewPassword ? 'Masquer' : 'Afficher'">
                    <svg v-if="!showResetNewPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                  </button>
                </div>
              </div>
              <div class="form-group">
                <label>{{ $t('auth.confirmPassword') }}</label>
                <div class="password-wrapper">
                  <input :type="showResetConfirmPassword ? 'text' : 'password'" v-model="confirmNewPassword" required placeholder="••••••••" />
                  <button type="button" class="toggle-password-btn" @click="showResetConfirmPassword = !showResetConfirmPassword" :title="showResetConfirmPassword ? 'Masquer' : 'Afficher'">
                    <svg v-if="!showResetConfirmPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                  </button>
                </div>
              </div>
              <button type="submit" :disabled="loading" class="btn btn-primary btn-block">
                {{ loading ? $t('auth.loading') : $t('auth.submitResetPassword') }}
              </button>
              <div style="text-align: center; margin-top: 1rem;">
                <a href="#" @click.prevent="forgotStep = 1" style="font-size: 0.85rem; color: #64748b; text-decoration: underline;">
                  {{ $t('auth.resendCode') }}
                </a>
              </div>
            </form>

            <form v-if="mode === 'reg-client'" @submit.prevent="handleRegisterClient" class="auth-form">
              <div class="form-row">
                <div class="form-group"><label>{{ $t('auth.firstname') }}</label><input v-model="regClient.firstname" required placeholder="Jean" /></div>
                <div class="form-group"><label>{{ $t('auth.lastname') }}</label><input v-model="regClient.lastname" required placeholder="Dupont" /></div>
              </div>
              <div class="form-group"><label>{{ $t('auth.email') }}</label><input v-model="regClient.email" type="email" required placeholder="jean.dupont@email.com" /></div>
              <div class="form-group"><label>{{ $t('auth.phone') }}</label><input v-model="regClient.phone" required placeholder="06 12 34 56 78" /></div>
              <div class="form-row">
                <div class="form-group">
                  <label>{{ $t('auth.password') }}</label>
                  <div class="password-wrapper">
                    <input :type="showClientPassword ? 'text' : 'password'" v-model="regClient.password" required placeholder="••••••••" />
                    <button type="button" class="toggle-password-btn" @click="showClientPassword = !showClientPassword" :title="showClientPassword ? 'Masquer' : 'Afficher'">
                      <svg v-if="!showClientPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                      <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                    </button>
                  </div>
                </div>
                <div class="form-group">
                  <label>{{ $t('auth.confirmPassword') }}</label>
                  <div class="password-wrapper">
                    <input :type="showClientConfirmPassword ? 'text' : 'password'" v-model="regClient.confirmPassword" required placeholder="••••••••" />
                    <button type="button" class="toggle-password-btn" @click="showClientConfirmPassword = !showClientConfirmPassword" :title="showClientConfirmPassword ? 'Masquer' : 'Afficher'">
                      <svg v-if="!showClientConfirmPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                      <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                    </button>
                  </div>
                </div>
              </div>
              <div class="form-group"><label>{{ $t('auth.address') }}</label><input v-model="regClient.address" required placeholder="12 Rue de la Paix, Paris" /></div>
              <button type="submit" :disabled="loading" class="btn btn-primary btn-block">
                {{ loading ? $t('auth.loading') : $t('auth.submitClient') }}
              </button>
            </form>

            <form v-if="mode === 'reg-merchant'" @submit.prevent="handleRegisterMerchant" class="auth-form">
              <div class="form-row">
                <div class="form-group"><label>{{ $t('auth.firstname') }}</label><input v-model="regMerchant.firstname" required placeholder="Marie" /></div>
                <div class="form-group"><label>{{ $t('auth.lastname') }}</label><input v-model="regMerchant.lastname" required placeholder="Curie" /></div>
              </div>
              <div class="form-row">
                <div class="form-group"><label>{{ $t('auth.email') }}</label><input v-model="regMerchant.email" type="email" required placeholder="contact@boulangerie.fr" /></div>
                <div class="form-group"><label>{{ $t('auth.phone') }}</label><input v-model="regMerchant.phone" required placeholder="01 23 45 67 89" /></div>
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label>{{ $t('auth.password') }}</label>
                  <div class="password-wrapper">
                    <input :type="showMerchantPassword ? 'text' : 'password'" v-model="regMerchant.password" required placeholder="••••••••" />
                    <button type="button" class="toggle-password-btn" @click="showMerchantPassword = !showMerchantPassword" :title="showMerchantPassword ? 'Masquer' : 'Afficher'">
                      <svg v-if="!showMerchantPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                      <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                    </button>
                  </div>
                </div>
                <div class="form-group">
                  <label>{{ $t('auth.confirmPassword') }}</label>
                  <div class="password-wrapper">
                    <input :type="showMerchantConfirmPassword ? 'text' : 'password'" v-model="regMerchant.confirmPassword" required placeholder="••••••••" />
                    <button type="button" class="toggle-password-btn" @click="showMerchantConfirmPassword = !showMerchantConfirmPassword" :title="showMerchantConfirmPassword ? 'Masquer' : 'Afficher'">
                      <svg v-if="!showMerchantConfirmPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                      <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                    </button>
                  </div>
                </div>
              </div>
              <div class="form-group"><label>{{ $t('auth.companyName') }}</label><input v-model="regMerchant.company_name" required placeholder="Boulangerie Bio" /></div>
              <div class="form-group"><label>{{ $t('auth.address') }}</label><input v-model="regMerchant.company_address" required placeholder="45 Avenue de la République, Paris" /></div>
              <div class="form-group"><label>{{ $t('auth.siret') }}</label><input v-model="regMerchant.siret_number" required placeholder="123 456 789 00012" /></div>
              <button type="submit" :disabled="loading" class="btn btn-primary btn-block">
                {{ loading ? $t('auth.loading') : $t('auth.submitMerchant') }}
              </button>
            </form>

            <div v-if="mode === 'contact-volunteer'" class="auth-form" style="text-align: center; padding: 1rem 0.2rem;">
              <h3 style="margin-bottom: 0.8rem; color: #2d6a4f; font-family: 'Outfit', sans-serif; font-size: 1.25rem; font-weight: 700;">
                {{ $t('auth.volunteerHeader') }}
              </h3>
              <p style="color: #64748b; font-size: 0.9rem; line-height: 1.6; margin-bottom: 1.5rem;">
                {{ $t('auth.volunteerBody') }}
              </p>
              <a href="mailto:nomorewasteesgi@gmail.com?subject=Candidature%20B%C3%A9n%C3%A9vole%20-%20No%20More%20Waste" class="btn btn-primary btn-block" style="display: inline-block; text-decoration: none; text-align: center; box-sizing: border-box; font-weight: 600;">
                {{ $t('auth.volunteerBtn') }}
              </a>
            </div>
          </div>
        </div>

      </div>
    </main>

    <footer class="landing-footer">
      <p>{{ $t('landing.footer') }}</p>
    </footer>
  </div>
</template>

<script>
import api from '../services/api'

export default {
  name: 'LoginView',
  data() {
    return {
      mode: 'login',
      email: '',
      password: '',
      showLoginPassword: false,
      showResetNewPassword: false,
      showResetConfirmPassword: false,
      showClientPassword: false,
      showClientConfirmPassword: false,
      showMerchantPassword: false,
      showMerchantConfirmPassword: false,
      loading: false,
      msg: '',
      msgType: '',
      regClient: { firstname: '', lastname: '', email: '', password: '', confirmPassword: '', address: '', phone: '' },
      regMerchant: { firstname: '', lastname: '', email: '', phone: '', password: '', confirmPassword: '', company_name: '', company_address: '', siret_number: '' },
      forgotStep: 1,
      forgotEmail: '',
      resetCode: '',
      newPassword: '',
      confirmNewPassword: ''
    }
  },
  methods: {
    flash(text, type) {
      this.msg = text
      this.msgType = type
      setTimeout(() => this.msg = '', 5000)
    },
    async handleLogin() {
      this.loading = true
      this.msg = ''
      try {
        const response = await api.post('/login', {
          email: this.email,
          password: this.password
        })
        const token = response.data.token
        localStorage.setItem('token', token)
        const decoded = JSON.parse(atob(token.split('.')[1]))
        const role = decoded.role
        localStorage.setItem('role', role)

        if (role === 'client') this.$router.push('/client/dashboard')
        else if (role === 'merchant') this.$router.push('/merchant/dashboard')
        else if (role === 'volunteer') this.$router.push('/volunteer/dashboard')
        else if (role === 'staff') this.$router.push('/staff/dashboard')
        else if (role === 'admin') this.$router.push('/admin/dashboard')
        else this.$router.push('/')
      } catch (error) {
        this.flash(error.response?.data?.error || 'Erreur de connexion au serveur', 'err')
      } finally {
        this.loading = false
      }
    },
    async handleForgotPassword() {
      this.loading = true
      this.msg = ''
      try {
        const res = await api.post('/forgot-password', { email: this.forgotEmail })
        this.flash(res.data.message || 'Si l\'adresse existe, un code vous a été envoyé.', 'ok')
        this.forgotStep = 2
      } catch (error) {
        this.flash(error.response?.data?.error || 'Erreur lors de la demande', 'err')
      } finally {
        this.loading = false
      }
    },
    async handleResetPassword() {
      if (this.newPassword !== this.confirmNewPassword) {
        this.flash('Les deux mots de passe ne correspondent pas', 'err')
        return
      }
      this.loading = true
      this.msg = ''
      try {
        const res = await api.post('/reset-password', {
          email: this.forgotEmail,
          code: this.resetCode,
          new_password: this.newPassword
        })
        this.flash(res.data.message || 'Mot de passe réinitialisé avec succès ! Connectez-vous.', 'ok')
        this.mode = 'login'
        this.email = this.forgotEmail
        this.password = ''
      } catch (error) {
        this.flash(error.response?.data?.error || 'Code invalide ou expiré', 'err')
      } finally {
        this.loading = false
      }
    },
    async handleRegisterClient() {
      if (this.regClient.password !== this.regClient.confirmPassword) {
        this.flash('Les deux mots de passe ne correspondent pas', 'err')
        return
      }
      this.loading = true
      this.msg = ''
      try {
        await api.post('/register/client', {
          firstname: this.regClient.firstname,
          lastname: this.regClient.lastname,
          email: this.regClient.email,
          password: this.regClient.password,
          address: this.regClient.address,
          phone: this.regClient.phone
        })
        this.flash('Compte client créé avec succès ! Connectez-vous.', 'ok')
        this.mode = 'login'
        this.email = this.regClient.email
        this.password = ''
      } catch (error) {
        this.flash(error.response?.data?.error || 'Erreur lors de l\'inscription', 'err')
      } finally {
        this.loading = false
      }
    },
    async handleRegisterMerchant() {
      if (this.regMerchant.password !== this.regMerchant.confirmPassword) {
        this.flash('Les deux mots de passe ne correspondent pas', 'err')
        return
      }
      this.loading = true
      this.msg = ''
      try {
        await api.post('/register/merchant', {
          firstname: this.regMerchant.firstname,
          lastname: this.regMerchant.lastname,
          email: this.regMerchant.email,
          phone: this.regMerchant.phone,
          password: this.regMerchant.password,
          company_name: this.regMerchant.company_name,
          company_address: this.regMerchant.company_address,
          siret_number: this.regMerchant.siret_number
        })
        this.flash('Compte commerçant soumis ! Il sera validé par notre équipe après vérification SIRET.', 'ok')
        this.mode = 'login'
        this.email = this.regMerchant.email
        this.password = ''
      } catch (error) {
        this.flash(error.response?.data?.error || 'Erreur lors de l\'inscription', 'err')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.landing-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0b1329 0%, #1e293b 50%, #0f172a 100%);
  color: #f8fafc;
  display: flex;
  flex-direction: column;
  font-family: 'Inter', sans-serif;
}

.hero-section {
  flex: 1;
  padding: 3rem 1.5rem;
  display: flex;
  align-items: center;
}

.hero-container {
  max-width: 1280px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 500px;
  gap: 3.5rem;
  align-items: center;
  width: 100%;
}

@media (max-width: 1024px) {
  .hero-container {
    grid-template-columns: 1fr;
    gap: 2.5rem;
  }
}

.hero-title {
  font-family: 'Outfit', sans-serif;
  font-size: 2.8rem;
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 1.2rem;
  letter-spacing: -0.02em;
  color: #ffffff;
}

.gradient-text {
  background: linear-gradient(135deg, #52b788, #74c69d);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-description {
  font-size: 1.05rem;
  color: #cbd5e1;
  line-height: 1.7;
  margin-bottom: 2.2rem;
  max-width: 620px;
}

.pillars-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2.5rem;
}

.pillar-card {
  background: rgba(30, 41, 59, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  padding: 1.1rem 1.3rem;
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  transition: all 0.25s ease;
}

.pillar-card:hover {
  background: rgba(30, 41, 59, 0.9);
  border-color: rgba(82, 183, 136, 0.4);
  transform: translateX(4px);
}

.pillar-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.pillar-icon-box.green { background: rgba(82, 183, 136, 0.2); }
.pillar-icon-box.orange { background: rgba(245, 158, 11, 0.2); }
.pillar-icon-box.blue { background: rgba(59, 130, 246, 0.2); }

.pillar-title {
  font-family: 'Outfit', sans-serif;
  font-size: 1.05rem;
  font-weight: 700;
  color: #f8fafc;
  margin-bottom: 0.2rem;
}

.pillar-desc {
  font-size: 0.85rem;
  color: #94a3b8;
  line-height: 1.45;
}

.impact-metrics {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 1rem 1.4rem;
  border-radius: 16px;
  width: fit-content;
}

.metric-item {
  display: flex;
  flex-direction: column;
}

.metric-value {
  font-family: 'Outfit', sans-serif;
  font-size: 1.5rem;
  font-weight: 800;
  color: #52b788;
  line-height: 1;
}

.metric-label {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 0.2rem;
  font-weight: 500;
}

.metric-divider {
  width: 1px;
  height: 28px;
  background: rgba(255, 255, 255, 0.12);
}

.hero-form-wrapper {
  width: 100%;
}

.login-card {
  background: #ffffff;
  padding: 2.2rem;
  border-radius: 24px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  width: 100%;
  box-sizing: border-box;
  color: #0f172a;
}

.login-header {
  text-align: center;
}

.login-header h2 {
  color: #0f172a;
  font-family: 'Outfit', 'Inter', sans-serif;
  font-size: 1.8rem;
  font-weight: 800;
  margin-bottom: 0.2rem;
  letter-spacing: -0.02em;
}

.login-subtitle {
  color: #64748b;
  font-size: 0.85rem;
  margin-bottom: 1.4rem;
}

.auth-tabs {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.35rem;
  margin-bottom: 1.5rem;
  background: #f1f5f9;
  padding: 5px;
  border-radius: 12px;
}

.auth-tabs button {
  padding: 0.55rem 0.2rem;
  border: none;
  background: transparent;
  font-family: 'Inter', sans-serif;
  font-size: 0.78rem;
  font-weight: 600;
  color: #64748b;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.auth-tabs button.active {
  background: #2d6a4f;
  color: #ffffff;
  box-shadow: 0 3px 10px rgba(45, 106, 79, 0.25);
  font-weight: 700;
}

.auth-tabs button:hover:not(.active) {
  background: #e2e8f0;
  color: #1e293b;
}

.auth-form {
  display: flex;
  flex-direction: column;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.9rem;
  width: 100%;
}

.form-group {
  display: flex;
  flex-direction: column;
  margin-bottom: 1rem;
  min-width: 0;
}

.form-group label {
  font-size: 0.82rem;
  font-weight: 600;
  color: #334155;
  margin-bottom: 0.3rem;
}

.form-group input {
  width: 100%;
  padding: 0.7rem 0.85rem;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  font-size: 0.88rem;
  box-sizing: border-box;
  font-family: 'Inter', sans-serif;
  transition: all 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #2d6a4f;
  box-shadow: 0 0 0 3px rgba(45, 106, 79, 0.15);
}

.password-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
}

.password-wrapper input {
  padding-right: 2.5rem;
}

.toggle-password-btn {
  position: absolute;
  right: 0.6rem;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  border-radius: 6px;
  transition: color 0.2s;
  z-index: 2;
}

.toggle-password-btn:hover {
  color: #2d6a4f;
}

.btn-block {
  width: 100%;
  margin-top: 0.6rem;
  padding: 0.8rem;
  font-size: 0.92rem;
  border-radius: 10px;
  font-weight: 600;
}

.landing-footer {
  text-align: center;
  padding: 1.5rem 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  font-size: 0.82rem;
  color: #64748b;
  background: rgba(15, 23, 42, 0.9);
}

@media (max-width: 640px) {
  .hero-title { font-size: 2.1rem; }
  .impact-metrics { flex-direction: column; align-items: flex-start; gap: 0.8rem; width: 100%; }
  .metric-divider { display: none; }
  .auth-tabs { grid-template-columns: repeat(2, 1fr); }
}
</style>