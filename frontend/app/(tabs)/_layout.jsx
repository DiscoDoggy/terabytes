import { Tabs } from 'expo-router';
import FontAwesome from "@expo/vector-icons/FontAwesome"

export default function TabLayout() {
    return (
        <Tabs screenOptions={{tabBarActiveTintyColor: "blue"}}>
            <Tabs.Screen 
                name = "index"
                options={{
                    title:"Feed",
                    tabBarIcon: ({ color }) => <FontAwesome size={28} name="home" color={color} />,
                }} 
            
            />
            <Tabs.Screen 
                name= "explore_screen"
                options={{
                    title: "Explore",
                    tabBarIcon:({ color }) => <FontAwesome size={28} name="search" color={color}/>,
                }} 
            />

            <Tabs.Screen 
                name= "create_post"
                options={{
                    title: "Post",
                    tabBarIcon:({ color }) => <FontAwesome size={28} name="plus-square" color={color}/>,
                }} 
            />

            <Tabs.Screen 
                name= "profile"
                options={{
                    title: "You",
                    tabBarIcon:({ color }) => <FontAwesome size={28} name="user" color={color}/>,
                }} 
            />
        </Tabs>
    );
}