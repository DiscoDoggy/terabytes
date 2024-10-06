import { Text, View, StyleSheet, Image } from "react-native";

const BlogFeedSummaryCard = (props) => {
    //expected props are title, author, company if applicable
    //an array of images, summary content, tags, and a path to read
    //the whole article

    return (
      <>
        
        <View style={styles.title_container}>
          <Text style={styles.title_text}>{props.title}</Text>
        </View>

        <View style={styles.authors_company_container}>
          <Text style={styles.authors_company_text}>{props.authors}
            <Text>{props.company}</Text>
          </Text>
        </View>

         <View style={styles.image_contanier}>
          <Image
            source={{
              uri: props.image_src
            }}
            style={{width:'100%', height: 200}}
          />
        </View> 

        <View>
          <Text>Article Description:</Text>
          <Text>{props.article_description}</Text>
        </View>
    </>
    );
};


const styles = StyleSheet.create({
  container: {
    flex: 1,
 
  },

  title_container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },

  title_text: {
    fontWeight: "bold",
    fontSize: 25,

  },

  authors_company_container: {
    flex: 1,
    flexDirection: 'row',
  },

  authors_company_text: {
    fontSize: 25,
    fontWeight: "medium",
  },

  image_container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center"
  }

});

export default BlogFeedSummaryCard;


