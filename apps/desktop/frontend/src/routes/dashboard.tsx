import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import RootLayout from '@/layouts/root-layout'
import { createFileRoute } from '@tanstack/react-router'
import { CalendarIcon, CreditCard, IndianRupee, Users } from 'lucide-react'

export const Route = createFileRoute('/dashboard')({
    component: () => (<RootLayout><DashboardPage /></RootLayout>)
})

const DashboardPage = () => {
    return (
        <main className="flex-1 overflow-auto p-4 md:p-6">
            <div className="grid gap-4 md:gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Sales</CardTitle>
                        <IndianRupee className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">₹5,231.89</div>
                        <p className="text-xs text-muted-foreground">
                            +20.1% from last month
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Tables</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">12</div>
                        <p className="text-xs text-muted-foreground">3 more than average</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Open Orders</CardTitle>
                        <CreditCard className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">23</div>
                        <p className="text-xs text-muted-foreground">5 pending payment</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Reservations</CardTitle>
                        <CalendarIcon className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">8</div>
                        <p className="text-xs text-muted-foreground">
                            2 upcoming in the next hour
                        </p>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="current-orders" className="mt-6">
                <TabsList className="grid w-full grid-cols-3">
                    <TabsTrigger value="current-orders">Current Orders</TabsTrigger>
                    <TabsTrigger value="menu">Menu</TabsTrigger>
                    <TabsTrigger value="tables">Tables</TabsTrigger>
                </TabsList>
                <TabsContent value="current-orders" className="mt-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Current Orders</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center">
                                    <span className="font-medium">Order #1234</span>
                                    <span className="text-sm text-muted-foreground">Table 5</span>
                                    <span className="text-sm font-semibold">₹45.00</span>
                                </div>
                                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center">
                                    <span className="font-medium">Order #1235</span>
                                    <span className="text-sm text-muted-foreground">Table 3</span>
                                    <span className="text-sm font-semibold">₹32.50</span>
                                </div>
                                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center">
                                    <span className="font-medium">Order #1236</span>
                                    <span className="text-sm text-muted-foreground">Takeout</span>
                                    <span className="text-sm font-semibold">₹27.75</span>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value="menu" className="mt-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Menu Items</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                <div className="flex justify-between items-center">
                                    <span className="font-medium">Margherita Pizza</span>
                                    <span className="text-sm font-semibold">$12.99</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="font-medium">Caesar Salad</span>
                                    <span className="text-sm font-semibold">$8.50</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="font-medium">Grilled Salmon</span>
                                    <span className="text-sm font-semibold">$18.75</span>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value="tables" className="mt-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Table Status</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                                <div className="p-4 bg-green-100 rounded-md text-center">
                                    <span className="font-medium">Table 1</span>
                                    <p className="text-sm text-green-600">Available</p>
                                </div>
                                <div className="p-4 bg-red-100 rounded-md text-center">
                                    <span className="font-medium">Table 2</span>
                                    <p className="text-sm text-red-600">Occupied</p>
                                </div>
                                <div className="p-4 bg-yellow-100 rounded-md text-center">
                                    <span className="font-medium">Table 3</span>
                                    <p className="text-sm text-yellow-600">Reserved</p>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </main>
    )
}
