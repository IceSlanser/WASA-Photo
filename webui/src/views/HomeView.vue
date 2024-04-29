<script>
import {RouterLink} from "vue-router";

export default {
  components: {RouterLink},
  data: function() {
    return {
      error: null,
      myID: localStorage.getItem("ID"),
      myUsername: localStorage.getItem("username"),
      newUsername: "",
      profileOwner: "",
      newText: "",

      userProfile: {
        profile: {
          ID: 0,
          username: "",
          following_count: 0,
          follower_count: 0,
          post_count: 0
        },
        posts: [
          {
            ID: 0,
            like_count: 0,
            comment_count: 0,
            showCommentInput: false
          }
        ],
        followings: [],
        followers: [],
        bannedFrom: [],
      },

      stream: [
        {
          post: {
            ID: 0,
            username: "",
            file: "",
            description: "",
            like_count: 0,
            comment_count: 0,
            date_time: ""
          },
          like_owners: [
            {
              username: ""
            }
          ],
          full_comments: [
            {
              username: "",
              comment: {
                ID: 0,
                post_ID: 0,
                owner_ID: 0,
                text: "",
                date_time:""
              }
            }
          ]
        }
      ],

      newComments: ["",],
      fullPost: {
        post: {
          ID: 0,
          like_count: 0,
          comment_count: 0,
        },
        like_owners: [
          {
            username: "",
          }
        ],
        full_comments: [
          {
            username: "",
            comment: {
              ID: 0,
              post_ID: 0,
              owner_ID: 0,
              owner_username: "",
              text: "",
              date_time: "",
            }
          }
        ],
      },

      showSearchInput: false,
      showCommentWindow: false,
      showLikeWindow: false,
      showLoading: false
    }
  },

  computed: {
    sortedPosts() {
      if (!this.stream) {
        this.stream = []
      }
      return this.stream.slice().sort((a, b) => new Date(b.post.date_time) - new Date(a.post.date_time));
    },
  },

  mounted() {
    this.getStream();
  },

	methods: {
    async doLogout() {
      await localStorage.clear()
      this.$router.push({path: '/'})
    },

    async getStream() {
      this.error = null
      this.showLoading = true;
      try {
        let response = await this.$axios.get(`/stream`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.stream = response.data
        await localStorage.setItem("stream", JSON.stringify(this.stream));
        this.showLoading = false;
      } catch (e) {
        this.showLoading = false;
        if (e.response && e.response.status === 400) {
          this.error = "Failed to get stream";
        } else if (e.response && e.response.status === 404) {
          this.error = "User not found";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later";
        } else {
          this.error = e.toString();
        }
      }
    },

    async getProfile() {
      await this.$router.push({path: `/users/${this.myID}/profile`})
    },

    async searchUser() {
      this.error = null
      try {
        let response = await this.$axios.get(`/users?username=${this.profileOwner}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.userProfile.profile.ID = response.data
        try {
          let res = await this.$axios.get(`/users/${this.userProfile.profile.ID}/profile`, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          })
          this.userProfile = res.data
          await localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
          this.$router.push({path: `/users/${this.userProfile.profile.ID}/profile`})
        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Failed to request user's profile.";
          } else if (e.response && e.response.status === 404) {
            this.error = "User not found.";
          } else if (e.response && e.response.status === 500) {
            this.error = "An internal error occurred, please try again later.";
          } else {
            this.error = e.toString();
          }
        }
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request user's profile.";
        } else if (e.response && e.response.status === 404) {
          this.error = "User not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async toggleSearchInput() {
      this.showSearchInput = !this.showSearchInput;
      if (!this.showSearchInput) {
        await localStorage.removeItem(this.newUsername)
        this.newUsername = ""
      }
    },

    async toggleLike(postID) {
      this.error = null;
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        const listOfLikeOwners = this.fullPost.like_owners || [];
        let isLiked = listOfLikeOwners.includes(this.myUsername);


        if (isLiked) {
          try {
            await this.$axios.delete(`/posts/${postID}/likes`, {
              headers: {
                Authorization: localStorage.getItem("username")
              }
            });
            this.userProfile.posts[i].like_count--;

          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Failed to delete.";
            } else if (e.response && e.response.status === 401) {
              this.error = "toggleLike not authorized";
            } else if (e.response && e.response.status === 404) {
              this.error = "Post not found";
            } else {
              this.error = e.toString();
            }
          }

        } else {
          try {
            await this.$axios.put(`/posts/${postID}/likes`, {}, {
              headers: {
                Authorization: localStorage.getItem("username")
              }
            });
            this.userProfile.posts[i].like_count++;

          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Failed to put.";
            } else if (e.response && e.response.status === 401) {
              this.error = "toggleLike not authorized";
            } else if (e.response && e.response.status === 404) {
              this.error = "Post not found";
            } else {
              this.error = e.toString();
            }
          }

        }
        await localStorage.setItem("userProfile", JSON.stringify(this.userProfile))
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post.";
        } else if (e.response && e.response.status === 401) {
          this.error = "toggleLike not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "Internal Server Error.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async showLikes(postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        await localStorage.setItem("fullPost", JSON.stringify(this.fullPost))
        this.showLikeWindow = true
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post.";
        } else if (e.response && e.response.status === 401) {
          this.error = "toggleLike not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "Internal Server Error.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async commentPhoto(postID) {
      try {
        let i = this.userProfile.posts.findIndex(post => post.ID === postID);
        let formData = new FormData();
        let tmp = this.newComments.reverse();
        formData.append('text', tmp[i])
        await this.$axios.post(`/posts/${postID}/comments`, formData, {
          headers: {
            Authorization: localStorage.getItem("username"),
          }
        })
        this.userProfile.posts[i].comment_count++;
        await localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        await this.toggleCommentInput(postID)
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request new comment.";
        } else if (e.response && e.response.status === 401) {
          this.error = "commentPhoto not authorized.";
        } else if (e.response && e.response.status === 404) {
          this.error = "Data not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async deleteComment(postID, commentID) {
      this.error = null
      try {
        await this.$axios.delete(`/posts/${postID}/comments/${commentID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        if (this.myUsername) {
          let i = this.userProfile.posts.findIndex(post => post.ID === postID);
          this.fullPost.full_comments = this.fullPost.full_comments.filter(comment => comment.post_ID !== comment.post_ID);
          this.userProfile.posts[i].comment_count--;
          await localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        }
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to delete.";
        } else if (e.response && e.response.status === 401) {
          this.error = "deleteComment not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Comment not found";
        } else {
          this.error = e.toString();
        }
      }
    },

    async showComments(postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        await localStorage.setItem("fullPost", JSON.stringify(this.fullPost))
        this.showCommentWindow = true
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post.";
        } else if (e.response && e.response.status === 401) {
          this.error = "toggleLike not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "Internal Server Error.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async toggleCommentInput(postID) {
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      this.userProfile.posts[i].showCommentInput = !this.userProfile.posts[i].showCommentInput;
      if (!this.userProfile.posts[i].showCommentInput) {
        this.newComments[i] = ""
      }
      await localStorage.setItem("userProfile", JSON.stringify(this.userProfile))
    },

    async closeLikeWindow() {
      this.showLikeWindow = false
    },

    async closeCommentWindow() {
      this.showCommentWindow = false
    },

  }
}

</script>

<template>
  <div>
    <div style="display: flex; justify-content: center;">
      <h2 class="h2">Home</h2>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>

    <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
  </div>


  <div class="loading-container" v-if="showLoading">
    <div class="loading">
      <h1 class="h1">Loading...</h1>
      <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#loader"/></svg>
    </div >
  </div>
  <div class="post-grid" v-else>
    <div v-for="(post, index) in sortedPosts" :key="post.ID" class="post-container" >
        <img v-if="post.post.file" :src="'data:image/jpeg;base64,' + post.post.file" alt="Post Image" class="post-image img-fluid align-content-center">
      <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
      <div class="d-flex justify-content-between">
        <p><span style="font-weight: bold;">{{ post.post.username }}</span>: {{ post.post.description }}</p>
        <p>{{ post.post.date_time }}</p>
      </div>
      <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
      <button type="button" class="btn" @click="showLikes(post.ID)">
        Likes: {{ post.post.like_count}}
      </button>
      <button type="button" class="btn mb-1" @click="toggleLike(post.ID)">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
      </button>
      <div>
        <button type="button" class="btn" @click="showComments(post.ID)">
          Comments: {{ post.post.comment_count }}
        </button>
        <button type="button" class="btn mb-1" @click="toggleCommentInput(post.ID)">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
        </button>
        <div v-if="post.post.showCommentInput" style="margin-right: 10px;">
          <input type="text" id="newComment" v-model="newComments[index]" class="form-control form-control-sm"
                 placeholder="comment text" aria-label="Recipient's comment" aria-describedby="basic-addon2">
          <button v-if="post.post.showCommentInput" type="button" class="btn btn-sm btn-primary" @click="commentPhoto(post.ID)">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
          </button>
        </div>
      </div>
    </div>
  </div>

  <div class="liked-users-overlay" v-if="this.showLikeWindow">
    <div class="liked-users-modal">
      <h2>Likes</h2>
      <ul>
        <li v-for="username in this.fullPost.like_owners" :key="username">{{ username }}</li>
      </ul>
      <button @click="this.closeLikeWindow()">Close</button>
    </div>
  </div>

  <div class="liked-users-overlay" v-if="this.showCommentWindow">
    <div class="liked-users-modal">
      <h2>Comments</h2>
      <ul>
        <li v-for="fullComment in this.fullPost.full_comments" :key="fullComment.username">
          {{ fullComment.username + ": " + fullComment.comment.text }}
          <button v-if="fullComment.username === this.myUsername" type="button" class="btn delete-button mb-2" @click="deleteComment(fullComment.comment.post_ID, fullComment.comment.ID)">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
          </button>
        </li>
      </ul>
      <button @click="this.closeCommentWindow()">Close</button>
    </div>
  </div>


  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky d-flex flex-column">
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/stream" class="nav-link" style="font-size: 20px;">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                Home
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink :to="'/users/' + myID + '/profile'" class="nav-link" style="font-size: 20px;" @click="getProfile">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                Profile
              </RouterLink>
            </li>
            <li class="nav-item mx-1">
              <button type="button" class="btn btn-sm" style="font-size: 20px;" @click="toggleSearchInput">
                <svg class="feather mx-1">
                  <use href="/feather-sprite-v4.29.0.svg#search"/>
                </svg>
                <span style="font-weight: 500;">Search</span>
              </button>
              <div class="d-flex ">
                <input v-if="showSearchInput" type="text" id="Searched Username" v-model="profileOwner"
                       class="form-control form-control-sm"
                       placeholder="Who are you searching for?" aria-label="Recipient's username" aria-describedby="basic-addon2">
                <button v-if="showSearchInput" type="button" class="btn btn-sm btn-primary ml-2 me-2" @click="searchUser">
                  Search
                </button>
              </div>
            </li>
          </ul>

          <div class="mt-auto mb-3 mx-1">
            <RouterLink to="/" class="nav-link" style="font-size: 20px;" @click="doLogout">
              <svg class="feather mx-1">
                <use href="/feather-sprite-v4.29.0.svg#log-out"/>
              </svg>
              Logout
            </RouterLink>
          </div>
        </div>
      </nav>
    </div>
  </div>

</template>

<style>
.post-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(calc(33.33% - 15px), 1fr));
  grid-gap: 15px;
}

.post-container {
  position: relative;
  border: 2px solid #ccc;
  margin-bottom: 20px;
  width: 100%;
  margin-top: 50px;
  justify-content: center;
}

.post-image {
  object-fit: cover;
  aspect-ratio: 1.91/1;
  display: block;
  margin: 0 auto;
  height: 275px;
}

.loading-container {
  top: 50%;
  left: 50%;
  width: 100%;
  height: 100%
}

.loading {
  text-align: center;
}

</style>
