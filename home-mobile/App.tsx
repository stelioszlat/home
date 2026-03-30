import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View } from 'react-native';
import * as Speech from 'expo-speech';

export default function App() {
   const speak = () => {
    const thingToSay = '1';
    Speech.speak(thingToSay);
  };
  return (
    <View style={styles.container}>
      <Button title="Press to hear some words" onPress={speak} />
    </View>
    <View style={styles.container}>
      <Text>Open up App.tsx to start working on your app!</Text>
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
