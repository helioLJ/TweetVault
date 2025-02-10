export interface Bookmark {
    id: string;
    created_at: string;
    full_text: string;
    screen_name: string;
    name: string;
    profile_image_url: string;
    favorite_count: number;
    retweet_count: number;
    bookmark_count: number;
    quote_count: number;
    reply_count: number;
    views_count: number;
    url: string;
    media: Media[];
    tags: Tag[];
    archived: boolean;
  }
  
  export interface Media {
    id: number;
    tweet_id: string;
    type: string;
    url: string;
    thumbnail: string;
    original: string;
    file_name: string;
  }
  
  export interface Tag {
    id: number;
    name: string;
    created_at: string;
  }