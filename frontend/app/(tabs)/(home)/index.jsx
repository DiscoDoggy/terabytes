import { Text, View, StyleSheet } from "react-native";

export default function HomeScreen() {
  return (
    <View styles={styles.container}>
      <Text>This is the home screen and will hold the feed</Text>
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