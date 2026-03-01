import { ofetch } from 'ofetch'

export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string

  const client = ofetch.create({
    baseURL
  })

  return {
    get: client,
    post: client,
    del: client
  }
}

