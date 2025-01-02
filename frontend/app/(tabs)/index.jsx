import React, { useEffect, useState } from "react";
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

  const [post_data, set_data] = useState([]);

  const get_followed_blog_posts = async () => {
    try {
      const response = await fetch("http://127.0.0.1:8000/posts/following", {
        headers: {
          "Content-Type" : "application/json",
          "Cookie" : "session_id=aba46209-e1e7-49d7-9ff2-e732d7bd0dad"

        }
      })
      if(!response.ok) {
        throw new Error(`Response Status: ${response.status}`);
      }
      
      const json = await response.json();
      set_data(json)
      console.log(json);
    } catch (error) {
      console.error(error.message);
    }
  };

  useEffect(() => {
    get_followed_blog_posts();
  }, []);
  
  return (
    <SafeAreaView styles={styles.container}>
      <FlatList
        data={post_data}
        renderItem= {({item}) => <BlogFeedSummaryCard  
          title={item.title}
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
