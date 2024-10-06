import { Text, View, StyleSheet } from "react-native";

export default function ExploreScreen() {
  return (
    <View styles={styles.container}>
        <Text>This is the Explore and it will hold 
        the exploration features
        </Text>
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