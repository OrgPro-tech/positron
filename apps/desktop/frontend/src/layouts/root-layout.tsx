import { Home, MenuIcon, ShoppingCart, Users, X } from 'lucide-react'
import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Link } from '@tanstack/react-router'

export default function RootLayout({
    children
}: {
    children: React.ReactElement
}) {
    const [sidebarOpen, setSidebarOpen] = useState(false)

    const [selectedOutlet, setSelectedOutlet] = useState("main")

    const outlets = [
        { id: "main", name: "Main Branch" },
        { id: "downtown", name: "Downtown" },
        { id: "airport", name: "Airport" },
    ]

    return (
        <div className="flex bg-gray-100 h-screen">
            {/* Mobile sidebar */}
            < div className={`fixed inset-0 z-50 bg-gray-800 bg-opacity-50 lg:hidden ${sidebarOpen ? "block" : "hidden"}`} onClick={() => setSidebarOpen(false)}></div >

            <div className={`fixed inset-y-0 left-0 z-50 w-64 bg-white transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0 ${sidebarOpen ? "translate-x-0" : "-translate-x-full"}`}>
                <div className="flex items-center justify-start py-2 px-4">
                    <h1 className="text-xl font-bold text-gray-800">Positron POS</h1>
                    <Button variant="ghost" size="icon" className="lg:hidden" onClick={() => setSidebarOpen(false)}>
                        <X className="h-6 w-6" />
                    </Button>
                </div>
                <nav className="mt-3">
                    <Link activeOptions={{
                        exact: true,
                    }} activeProps={{
                        className: "bg-gray-200"
                    }} to="/dashboard" className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-200">
                        <Home className="h-5 w-5 mr-3" />
                        Dashboard
                    </Link>
                    <Link activeOptions={{
                        exact: true
                    }} activeProps={{
                        className: "bg-gray-200"
                    }} to="/menus" className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-200">
                        <MenuIcon className="h-5 w-5 mr-3" />
                        Menu
                    </Link>
                    <Link href="/orders" className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-200">
                        <ShoppingCart className="h-5 w-5 mr-3" />
                        Orders
                    </Link>
                    <Link href="/customers" className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-200">
                        <Users className="h-5 w-5 mr-3" />
                        Customers
                    </Link>
                </nav>
            </div>

            <div className="flex-1 flex flex-col overflow-hidden">
                {/* Header */}
                <header className="bg-white shadow-sm py-2 flex justify-end items-center">
                    <div className="pr-4">
                        <Select value={selectedOutlet} onValueChange={setSelectedOutlet}>
                            <SelectTrigger className="w-[180px]">
                                <SelectValue placeholder="Select outlet" />
                            </SelectTrigger>
                            <SelectContent>
                                {outlets.map((outlet) => (
                                    <SelectItem key={outlet.id} value={outlet.id}>
                                        {outlet.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                </header>
                <div className="overflow-visible overflow-y-scroll">
                    {children}
                </div>
            </div>
        </div>
    )
}
