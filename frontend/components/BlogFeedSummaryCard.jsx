import { Text, View, StyleSheet, Image } from "react-native";

const BlogFeedSummaryCard = (props) => {
    //expected props are title, author, company if applicable

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

        <View style={styles.article_summary_container}>
          <Text style={{fontWeight:"bold"}}>Summary:</Text>
          <Text style={styles.article_summary_text}>{props.article_description}</Text>
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
    marginLeft: 5
  },

  authors_company_text: {
    fontSize: 12,
    fontWeight: "light",

  },

  article_summary_container: {
    marginLeft: 5
  },

  article_summary_text: {
    fontSize: 15,
  },

  image_container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center"
  }

});

export default BlogFeedSummaryCard;


