import { createRootRoute, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'

export const Route = createRootRoute({
    component: () => <MainRoot />
})

function MainRoot() {
    return (
        <>
            <Outlet />
            <TanStackRouterDevtools />
        </ >
    )
}
