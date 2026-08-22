<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h2>{{ $t('auth.title') }}</h2>
        <p class="login-subtitle">{{ $t('auth.subtitle') }}</p>
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
        <div style="font-size: 2.5rem; margin-bottom: 0.5rem;"></div>
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
      if (!this.forgotEmail) return
      this.loading = true
      try {
        const res = await api.post('/forgot-password', { email: this.forgotEmail })
        this.flash(res.data.message || 'Code envoyé ! (Code de démonstration: 123456)', 'ok')
        this.resetCode = res.data.code || '123456'
        this.forgotStep = 2
      } catch (error) {
        this.flash(error.response?.data?.error || 'Erreur lors de la demande', 'err')
      } finally {
        this.loading = false
      }
    },
    async handleResetPassword() {
      if (!this.resetCode || !this.newPassword) return
      if (this.newPassword !== this.confirmNewPassword) {
        this.flash(this.$t('auth.passwordMismatch'), 'err')
        return
      }
      this.loading = true
      try {
        const res = await api.post('/reset-password', {
          email: this.forgotEmail,
          code: this.resetCode,
          new_password: this.newPassword
        })
        this.flash(res.data.message || 'Mot de passe réinitialisé avec succès !', 'ok')
        this.email = this.forgotEmail
        this.password = this.newPassword
        this.mode = 'login'
        this.forgotStep = 1
        this.resetCode = ''
        this.newPassword = ''
        this.confirmNewPassword = ''
      } catch (error) {
        this.flash(error.response?.data?.error || 'Erreur de réinitialisation', 'err')
      } finally {
        this.loading = false
      }
    },
    async handleRegisterClient() {
      if (this.regClient.password !== this.regClient.confirmPassword) {
        this.flash(this.$t('auth.passwordMismatch'), 'err')
        return
      }
      this.loading = true
      try {
        const payload = {
          firstname: this.regClient.firstname,
          lastname: this.regClient.lastname,
          email: this.regClient.email,
          password: this.regClient.password,
          address: this.regClient.address
        }
        await api.post('/register/client', payload)
        this.flash('Compte client créé avec succès ! Vous pouvez maintenant vous connecter.', 'ok')
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
        this.flash(this.$t('auth.passwordMismatch'), 'err')
        return
      }
      this.loading = true
      try {
        const payload = {
          firstname: this.regMerchant.firstname,
          last_name: this.regMerchant.lastname,
          email: this.regMerchant.email,
          phone: this.regMerchant.phone,
          password: this.regMerchant.password,
          company_name: this.regMerchant.company_name,
          company_address: this.regMerchant.company_address,
          siret_number: this.regMerchant.siret_number,
          address: this.regMerchant.company_address
        }
        await api.post('/register/merchant', payload)
        this.flash('Demande commerçant soumise avec succès ! En attente d\'approbation par l\'équipe staff.', 'ok')
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
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  padding: 1.5rem 1rem;
  box-sizing: border-box;
}

.login-card {
  background: #ffffff;
  padding: 2.2rem;
  border-radius: 20px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.4);
  width: 100%;
  max-width: 540px;
  box-sizing: border-box;
  transition: all 0.3s ease;
}

.login-header {
  text-align: center;
}

.login-header h2 {
  color: #0f172a;
  font-family: 'Outfit', 'Inter', sans-serif;
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 0.2rem;
  letter-spacing: -0.02em;
}

.login-subtitle {
  color: #64748b;
  font-size: 0.88rem;
  margin-bottom: 1.4rem;
}

/* TABS ADAPTATIFS & RESPONSIVES */
.auth-tabs {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
  margin-bottom: 1.6rem;
  background: #f8fafc;
  padding: 6px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
}

@media (min-width: 600px) {
  .auth-tabs {
    grid-template-columns: repeat(4, 1fr);
  }
}

.auth-tabs button {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 0.6rem 0.3rem;
  border: none;
  background: transparent;
  font-family: 'Inter', sans-serif;
  font-size: 0.8rem;
  font-weight: 600;
  color: #64748b;
  border-radius: 9px;
  cursor: pointer;
  transition: all 0.2s ease;
  line-height: 1.25;
  min-height: 44px;
  box-sizing: border-box;
}

.auth-tabs button.active {
  background: #2d6a4f;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(45, 106, 79, 0.3);
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
  gap: 1rem;
  width: 100%;
}

@media (max-width: 480px) {
  .login-card {
    padding: 1.5rem 1rem;
  }
  .form-row {
    grid-template-columns: 1fr;
  }
  .auth-tabs {
    grid-template-columns: repeat(2, 1fr);
  }
}

.form-group {
  display: flex;
  flex-direction: column;
  margin-bottom: 1.1rem;
  min-width: 0;
}

.form-group label {
  font-size: 0.83rem;
  font-weight: 600;
  color: #334155;
  margin-bottom: 0.35rem;
}

.form-group input {
  width: 100%;
  padding: 0.75rem 0.9rem;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  font-size: 0.9rem;
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
  padding-right: 2.6rem;
}

.toggle-password-btn {
  position: absolute;
  right: 0.6rem;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.35rem;
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
  margin-top: 0.8rem;
  padding: 0.85rem;
  font-size: 0.95rem;
  border-radius: 10px;
  font-weight: 600;
}
</style>