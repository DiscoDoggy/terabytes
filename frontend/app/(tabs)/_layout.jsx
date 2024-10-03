import { Tabs } from 'expo-router';
import FontAwesome from "@expo/vector-icons/FontAwesome"

export default function TabLayout() {
    return (
        <Tabs screenOptions={{tabBarActiveTintyColor: "blue"}}>
            <Tabs.Screen 
                name = "(home)"
                options={{
                    title:"Feed",
                    tabBarIcon: ({ color }) => <FontAwesome size={28} name="home" color={color} />,
                }} 
            
            />
            <Tabs.Screen 
                name= "(explore)"
                options={{
                    title: "Explore",
                    tabBarIcon:({ color }) => <FontAwesome size={28} name="search" color={color}/>,
                }} 
            />
        </Tabs>
    );
}