import { createFileRoute } from '@tanstack/react-router'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { Leaf, Flame, Check, X } from 'lucide-react'
import RootLayout from '@/layouts/root-layout'

export const Route = createFileRoute('/menus')({
    component: () => {
        return <RootLayout><MenuPage /></RootLayout>
    }
})

// Define the MenuItem type based on the provided structure
type MenuItem = {
    category: string
    name: string
    description: string
    price: number
    is_vegetarian: boolean
    spice_level: 'Mild' | 'Medium' | 'Hot'
    is_available: boolean
    image: string
}

// Sample menu items
const menuItems: MenuItem[] = [
    {
        category: 'Starter',
        name: 'Margherita Pizza',
        description: 'Classic Italian pizza with tomato and mozzarella',
        price: 250,
        is_vegetarian: true,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://images.prismic.io/eataly-us/ed3fcec7-7994-426d-a5e4-a24be5a95afd_pizza-recipe-main.jpg?auto=compress,format',
    },
    {
        category: 'Main Course',
        name: 'Spicy Chicken Curry',
        description: 'Tender chicken pieces in a spicy curry sauce',
        price: 304,
        is_vegetarian: false,
        spice_level: 'Hot',
        is_available: true,
        image: 'https://www.foodandwine.com/thmb/8YAIANQTZnGpVWj2XgY0dYH1V4I=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/spicy-chicken-curry-FT-RECIPE0321-58f84fdf7b484e7f86894203eb7834e7.jpg',
    },
    {
        category: 'Dessert',
        name: 'Chocolate Lava Cake',
        description: 'Warm chocolate cake with a gooey center',
        price: 149,
        is_vegetarian: true,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://www.foodandwine.com/thmb/XdFd-DvTtouryLCjeCqwhfmmK-A=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/molten-chocolate-cake-FT-RECIPE0220-0a33d7d0ab0c45588f7bfe742d33a9bc.jpg',
    },
    {
        category: 'Starter',
        name: 'Bruschetta',
        description: 'Grilled bread rubbed with garlic and topped with diced tomatoes and fresh basil',
        price: 180,
        is_vegetarian: true,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://www.allrecipes.com/thmb/QSsjryxShEx1L6o0HLer1Nn4jwA=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/54165-balsamic-bruschetta-DDMFS-4x3-e2b55b5ca39b4c1783e524a2461634ea.jpg',
    },
    {
        category: 'Main Course',
        name: 'Grilled Salmon',
        description: 'Fresh Atlantic salmon fillet grilled to perfection, served with lemon butter sauce',
        price: 350,
        is_vegetarian: false,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://www.pccmarkets.com/wp-content/uploads/2017/08/pcc-rosemary-grilled-salmon-flo.jpg',
    },
    {
        category: 'Main Course',
        name: 'Vegetable Stir Fry',
        description: 'Assorted fresh vegetables stir-fried in a savory sauce, served with steamed rice',
        price: 280,
        is_vegetarian: true,
        spice_level: 'Medium',
        is_available: true,
        image: 'https://images.immediate.co.uk/production/volatile/sites/30/2022/10/vegetarian-stir-fry-hero-fe84012.jpg?quality=90&resize=556,505',
    },
    {
        category: 'Starter',
        name: 'Spinach and Artichoke Dip',
        description: 'Creamy dip with spinach, artichoke hearts, and melted cheese, served with tortilla chips',
        price: 220,
        is_vegetarian: true,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://www.thespruceeats.com/thmb/IzI21XI-Gg07LQnFEu57xYVnA7w=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/spinach-and-artichoke-dip-4157518-07-8685813570f34ac89c090084c042266d.jpg',
    },
    {
        category: 'Main Course',
        name: 'Beef Tenderloin Steak',
        description: '8oz beef tenderloin cooked to your liking, served with mashed potatoes and seasonal vegetables',
        price: 450,
        is_vegetarian: false,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://i2.wp.com/www.downshiftology.com/wp-content/uploads/2023/02/How-To-Make-Filet-Mignon-5-600x400.jpg',
    },
    {
        category: 'Dessert',
        name: 'New York Cheesecake',
        description: 'Rich and creamy cheesecake with a graham cracker crust, topped with fresh berries',
        price: 180,
        is_vegetarian: true,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://www.allrecipes.com/thmb/v8JZdIICA1oerzX0L-KzyW2w9hM=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/8350-chantals-new-york-cheesecake-DDMFS-4x3-426569e82b4142a6a1ed01e068544245.jpg',
    },
    {
        category: 'Starter',
        name: 'Crispy Calamari',
        description: 'Lightly battered and fried calamari rings, served with marinara sauce',
        price: 240,
        is_vegetarian: false,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://i0.wp.com/www.russianfilipinokitchen.com/wp-content/uploads/2015/04/crispy-fried-calamari-01.jpg',
    },
    {
        category: 'Main Course',
        name: 'Vegetarian Lasagna',
        description: 'Layers of pasta, ricotta cheese, and vegetables, baked with marinara sauce and mozzarella',
        price: 320,
        is_vegetarian: true,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://myfoodstory.com/wp-content/uploads/2019/11/Easy-Vegetarian-Lasagna-5.jpg',
    },
    {
        category: 'Dessert',
        name: 'Tiramisu',
        description: 'Classic Italian dessert with layers of coffee-soaked ladyfingers and mascarpone cream',
        price: 160,
        is_vegetarian: true,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://images.squarespace-cdn.com/content/v1/5eed333596cab776eee55b17/e75e6aea-2e2a-4e5e-ad56-45f8d0a7fd63/AdobeStock_273554640.jpeg',
    },
    {
        category: 'Main Course',
        name: 'Shrimp Scampi',
        description: 'Succulent shrimp sautéed in garlic butter sauce, served over linguine pasta',
        price: 380,
        is_vegetarian: false,
        spice_level: 'Mild',
        is_available: true,
        image: 'https://www.allrecipes.com/thmb/jiV_4f8vXFle1RdFLgd8-_31J3M=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/229960-shrimp-scampi-with-pasta-DDMFS-4x3-e065ddef4e6d44479d37b4523808cc23.jpg',
    }
]

export default function MenuPage() {
    const [searchTerm, setSearchTerm] = useState('')
    const [categoryFilter, setCategoryFilter] = useState('All')

    const categories = ['All', ...new Set(menuItems.map((item) => item.category))]

    const filteredItems = menuItems.filter(
        (item) =>
            (categoryFilter === 'All' || item.category === categoryFilter) &&
            item.name.toLowerCase().includes(searchTerm.toLowerCase()),
    )

    return (
        <div className="p-4">
            <h1 className="text-3xl font-bold mb-6">Our Menu</h1>
            <div className="flex flex-col md:flex-row justify-between mb-6 space-y-4 md:space-y-0 md:space-x-4">
                <Input
                    type="search"
                    placeholder="Search menu items..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="max-w-sm"
                />
                <Select value={categoryFilter} onValueChange={setCategoryFilter}>
                    <SelectTrigger className="max-w-[180px]">
                        <SelectValue placeholder="Select category" />
                    </SelectTrigger>
                    <SelectContent>
                        {categories.map((category) => (
                            <SelectItem key={category} value={category}>
                                {category}
                            </SelectItem>
                        ))}
                    </SelectContent>
                </Select>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {filteredItems.map((item) => (
                    <Card key={item.name} className="flex flex-col">
                        <CardHeader>
                            <div className="relative h-48 w-full mb-4">
                                <img
                                    src={item.image}
                                    alt={item.name}
                                    className="absolute inset-0 w-full h-full object-cover rounded-t-lg"
                                />
                            </div>
                            <CardTitle className="flex justify-between items-center">
                                <span>{item.name}</span>
                                <span className="text-lg font-normal">
                                    ₹{item.price.toFixed(2)}
                                </span>
                            </CardTitle>
                            <CardDescription>{item.description}</CardDescription>
                        </CardHeader>
                        <CardContent className="flex-grow">
                            <div className="flex items-center space-x-2 mb-2">
                                {item.is_vegetarian && (
                                    <Leaf className="text-green-500" size={20} />
                                )}
                                {item.spice_level !== 'Mild' && (
                                    <Flame
                                        className={
                                            item.spice_level === 'Hot'
                                                ? 'text-red-500'
                                                : 'text-orange-400'
                                        }
                                        size={20}
                                    />
                                )}
                            </div>
                        </CardContent>
                        <CardFooter className="flex justify-between items-center">
                            <span className="text-sm font-medium">
                                {item.is_available ? (
                                    <span className="text-green-600 flex items-center">
                                        <Check size={16} className="mr-1" /> Available
                                    </span>
                                ) : (
                                    <span className="text-red-600 flex items-center">
                                        <X size={16} className="mr-1" /> Unavailable
                                    </span>
                                )}
                            </span>
                            <Button disabled={!item.is_available}>
                                {item.is_available ? 'Add to Order' : 'Out of Stock'}
                            </Button>
                        </CardFooter>
                    </Card>
                ))}
            </div>
        </div>
    )
}
