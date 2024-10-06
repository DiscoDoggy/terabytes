import { Text, View, StyleSheet } from "react-native";

export default function HomeScreen() {
  return (
    <View styles={styles.container}>
      <Text>This is the My Profile Screen and will hold the account info of the user</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: 
  {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  }
});