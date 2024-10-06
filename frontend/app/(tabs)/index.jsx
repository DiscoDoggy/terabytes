import { Text, View, StyleSheet, SafeAreaView, FlatList } from "react-native";
import BlogFeedSummaryCard from "@/components/BlogFeedSummaryCard";

export default function HomeScreen() {

  const POSTS_DATA = [
    {
      id:"1",
      title:"Assess, Test and Prepare for the (Un)Expected News Event",
      authors: "NYT Open Team: ",
      company: "NyTimes",
      image_src: "https://miro.medium.com/v2/resize:fit:1100/format:webp/1*u4RAa6kFb2pq-0ioKZ1Cdg.jpeg",
      article_description: "A brief history of New York Times Engineerings election readiness efforts’ to achieve reliable and resilient systems."
    },

    {
      id:"2",
      title:"Assess, Test and Prepare for the (Un)Expected News Event",
      authors: "NYT Open Team: ",
      company: "NyTimes",
      image_src: "https://miro.medium.com/v2/resize:fit:1100/format:webp/1*u4RAa6kFb2pq-0ioKZ1Cdg.jpeg",
      article_description: "A brief history of New York Times Engineerings election readiness efforts’ to achieve reliable and resilient systems."
    },
    
    {
      id:"3",
      title:"Assess, Test and Prepare for the (Un)Expected News Event",
      authors: "NYT Open Team: ",
      company: "NyTimes",
      image_src: "https://miro.medium.com/v2/resize:fit:1100/format:webp/1*u4RAa6kFb2pq-0ioKZ1Cdg.jpeg",
      article_description: "A brief history of New York Times Engineerings election readiness efforts’ to achieve reliable and resilient systems."
    },

    {
      id:"4",
      title:"Assess, Test and Prepare for the (Un)Expected News Event",
      authors: "NYT Open Team: ",
      company: "NyTimes",
      image_src: "https://miro.medium.com/v2/resize:fit:1100/format:webp/1*u4RAa6kFb2pq-0ioKZ1Cdg.jpeg",
      article_description: "A brief history of New York Times Engineerings election readiness efforts’ to achieve reliable and resilient systems."
    }

  ]
  
  return (
    <SafeAreaView styles={styles.container}>
      <FlatList
        data={POSTS_DATA}
        renderItem= {({item}) => <BlogFeedSummaryCard  
          title={item.title}
          authors={item.authors}
          company={item.company}
          image_src={item.image_src}
          article_description={item.article_description}
        />}
        keyExtractor={item=>item.id}
      />
    </SafeAreaView>

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
