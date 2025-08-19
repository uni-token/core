<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/composables/auth'

const props = defineProps<{
  register?: boolean
}>()

const { t } = useI18n()
const { login, register, isLoading } = useAuth()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const errorMessage = ref('')

async function handleSubmit() {
  if (!username.value || !password.value) {
    errorMessage.value = t('loginForm.fillCompleteInfo')
    return
  }

  if (props.register) {
    if (!confirmPassword.value) {
      errorMessage.value = t('loginForm.confirmPasswordRequired')
      return
    }
    if (password.value !== confirmPassword.value) {
      errorMessage.value = t('loginForm.passwordMismatch')
      return
    }
  }

  errorMessage.value = ''

  try {
    const result = props.register
      ? await register(username.value, password.value)
      : await login(username.value, password.value)

    if (result.status !== 'success') {
      errorMessage.value = result.message || t('loginForm.operationFailed')
    }
  }
  catch (error) {
    errorMessage.value = `${t('loginForm.networkError')}: ${error instanceof Error ? error.message : t('loginForm.unknownError')}`
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background p-4">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <CardTitle class="text-2xl font-bold">
          {{ props.register ? t('loginForm.register') : t('loginForm.login') }}
        </CardTitle>
        <CardDescription>
          {{ props.register ? t('loginForm.createAccount') : t('loginForm.enterCredentials') }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          class="space-y-4"
          :aria-label="props.register ? t('loginForm.registerForm') : t('loginForm.loginForm')"
          novalidate
          @submit.prevent="handleSubmit"
        >
          <div class="space-y-2">
            <label for="username" class="text-sm font-medium">{{ t('loginForm.username') }}</label>
            <Input
              id="username"
              v-model="username"
              type="text"
              :placeholder="t('loginForm.usernamePlaceholder')"
              required
              :disabled="isLoading"
              autocomplete="username"
              :aria-invalid="errorMessage ? 'true' : 'false'"
              :aria-describedby="errorMessage ? 'error-message' : undefined"
            />
          </div>
          <div class="space-y-2">
            <label for="password" class="text-sm font-medium">{{ t('loginForm.password') }}</label>
            <Input
              id="password"
              v-model="password"
              type="password"
              :placeholder="t('loginForm.passwordPlaceholder')"
              required
              :disabled="isLoading"
              :autocomplete="props.register ? 'new-password' : 'current-password'"
              :aria-invalid="errorMessage ? 'true' : 'false'"
              :aria-describedby="errorMessage ? 'error-message' : undefined"
            />
          </div>
          <div v-if="props.register" class="space-y-2">
            <label for="confirmPassword" class="text-sm font-medium">{{ t('loginForm.confirmPassword') }}</label>
            <Input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              :placeholder="t('loginForm.confirmPasswordPlaceholder')"
              required
              :disabled="isLoading"
              autocomplete="new-password"
              :aria-invalid="errorMessage ? 'true' : 'false'"
              :aria-describedby="errorMessage ? 'error-message' : undefined"
              :aria-label="t('loginForm.confirmPassword')"
            />
          </div>
          <Button
            type="submit"
            class="w-full"
            :disabled="isLoading"
            :aria-label="isLoading ? t('loginForm.processingWait') : (props.register ? t('loginForm.submitRegister') : t('loginForm.submitLogin'))"
          >
            <div v-if="isLoading" class="flex items-center gap-2">
              <div
                class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"
                aria-hidden="true"
              />
              {{ t('loginForm.processing') }}
            </div>
            <span v-else>{{ props.register ? t('loginForm.register') : t('loginForm.login') }}</span>
          </Button>
          <div
            v-if="errorMessage"
            id="error-message"
            class="text-sm text-destructive text-center"
            role="alert"
            aria-live="polite"
          >
            {{ errorMessage }}
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
