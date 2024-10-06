import { Text, View, StyleSheet } from "react-native";

export default function HomeScreen() {
  return (
    <View styles={styles.container}>
      <Text>This is the create post screen and will hold the text editor to create a post</Text>
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