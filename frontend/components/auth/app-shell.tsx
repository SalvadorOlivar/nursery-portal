'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect, type ReactNode } from 'react'
import { Button } from '@/components/ui/button'
import { useLogout, useMe } from '@/features/auth/hooks/use-auth'

const publicPaths = ['/login']

export function AppShell({ children }: { children: ReactNode }) {
  const pathname = usePathname()
  const router = useRouter()
  const isPublic = publicPaths.includes(pathname)
  const { data, isLoading, isError } = useMe()
  const logoutMutation = useLogout()
  const user = data?.user
  const canViewEmployees = user?.role === 'ADMIN' || user?.role === 'SUPERVISOR'

  useEffect(() => {
    if (!isPublic && isError) {
      router.replace('/login')
    }
  }, [isError, isPublic, router])

  async function handleLogout() {
    await logoutMutation.mutateAsync()
    router.replace('/login')
  }

  if (!isPublic && isLoading) {
    return (
      <div className="min-h-screen grid place-items-center text-sm text-muted-foreground">
        Cargando sesion...
      </div>
    )
  }

  if (isPublic) {
    return (
      <main className="min-h-screen">
        {children}
      </main>
    )
  }

  if (!user) {
    return null
  }

  return (
    <>
      <header className="border-b">
        <div className="container mx-auto px-4 flex items-center justify-between h-14">
          <Link href="/planificaciones" className="font-semibold text-lg">
            Nursery Portal
          </Link>
          <nav className="flex items-center gap-4 text-sm">
            {canViewEmployees && (
              <Link href="/employees" className="hover:text-primary transition-colors">
                Empleados
              </Link>
            )}
            <Link href="/intercambio" className="hover:text-primary transition-colors">
              Intercambios
            </Link>
            <Link href="/leave-requests" className="hover:text-primary transition-colors">
              Licencias
            </Link>
            <Link href="/planificaciones" className="hover:text-primary transition-colors">
              Planificaciones
            </Link>
            <Link href="/profile" className="text-muted-foreground hover:text-primary transition-colors">
              {user.username}
            </Link>
            <Button variant="outline" size="sm" onClick={handleLogout} disabled={logoutMutation.isPending}>
              Salir
            </Button>
          </nav>
        </div>
      </header>
      <main className="flex-1 container mx-auto px-4 py-6">
        {children}
      </main>
    </>
  )
}
