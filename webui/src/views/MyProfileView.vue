<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  components: {RouterLink},
  data: function () {
    return {
      error: null,
      userProfile: JSON.parse(localStorage.getItem("userProfile")),
      username: localStorage.getItem("username"),
      newUsername: "",
      newPhoto: "",
      newText: "",
      fullPost: {
        post: {
          ID: 0,
          like_count: 0,
          comment_count: 0,
        },
        like_owners: [
          {
            username: 0
          }
        ],
        full_comments: [
          {
            username: "",
            comment: {
              ID: 0,
              owner_ID: 0,
              owner_username: "",
              text: "",
            }
          }
        ]
      },
      showUsernameInput: false,
      showUploadInput: false,
      showCommentInput: false,
      showLikeWindow: false,
      showCommentWindow: false
    }
  },
  methods: {
    async doLogout() {
      localStorage.clear()
      this.$router.push({path: '/'})
    },
    async getMyProfile() {
      this.error = null
      try {
        let response = await this.$axios.get(`/users/${this.userProfile.profile.ID}/profile`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.userProfile = response.data
        localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
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
    },

    handleFileChange(event) {
      this.newPhoto = event.target.files[0];
      const reader = new FileReader();

      reader.readAsArrayBuffer(this.newPhoto);
    },
    async uploadPhoto() {
      try {
        let formData = new FormData();
        formData.append('file', this.newPhoto);
        formData.append('description', this.newText)
        await this.$axios.post("/posts", formData, {
          headers: {
            Authorization: localStorage.getItem("username"),
            "Content-Type": "multipart/form-data"
          }
        })
        await this.getMyProfile();
        window.location.reload()
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request new username.";
        } else if (e.response && e.response.status === 401) {
          this.error = "uploadPhoto not authorized.";
        } else if (e.response && e.response.status === 404) {
          this.error = "Data not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async deletePhoto(postID) {
      this.error = null
      try {
        await this.$axios.delete(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        if (this.username) {
          this.userProfile.posts = this.userProfile.posts.filter(post => post.ID !== postID);
          localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        }
        window.location.reload()
      } catch (e){
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request new username.";
        } else if (e.response && e.response.status === 401) {
          this.error = "setMyUserName not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found";
        } else {
          this.error = e.toString();
        }
      }
    },

    async setMyUserName() {
      if (this.username === "") {
        this.error = "Username cannot be empty.";
      } else {
        this.error = null
        try {
          await this.$axios.put("/profile/setUserName", {username: this.newUsername}, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          })
          localStorage.setItem("username", this.newUsername);
          this.userProfile.profile.username = this.newUsername;
          localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
          window.location.reload()
        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Failed to request new username.";
          } else if (e.response && e.response.status === 401) {
            this.error = "setMyUserName not authorized";
          } else if (e.response && e.response.status === 500) {
            this.error = "An internal error occurred, please try again later.";
          } else {
            this.error = e.toString();
          }
        }
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
        const likeOwnersArray = this.fullPost.like_owners || [];
        let isLiked = likeOwnersArray.includes(this.userProfile.profile.ID);

        if (isLiked) {
          await this.$axios.delete(`/posts/${postID}/likes`, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.posts[i].like_count--;

        } else {
          await this.$axios.put(`/posts/${postID}/likes`, {}, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.posts[i].like_count++;
        }
        localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        window.location.reload()
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

    async showLikes (postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
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
        let formData = new FormData();
        formData.append('text', this.newText)
        await this.$axios.post(`/posts/${postID}/comments`, formData, {
          headers: {
            Authorization: localStorage.getItem("username"),
          }
        })
        await this.getMyProfile();
        window.location.reload()
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

    async showComments (postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
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

    async closeLikeWindow () {
      this.showLikeWindow = false
    },

    async closeCommentWindow () {
      this.showCommentWindow = false
    },

    async toggleUsernameInput() {
      this.showUsernameInput = !this.showUsernameInput;
      if (!this.showUsernameInput) {
        localStorage.removeItem(this.newUsername)
        this.newUsername = ""
      }
    },
    async togglePhotoInput() {
      this.showUploadInput = !this.showUploadInput;
      if (!this.showUploadInput) {
        localStorage.removeItem(this.newPhoto)
        this.newPhoto = ""
        this.newText = ""
      }
    },
    async toggleCommentInput() {
      this.showCommentInput = !this.showCommentInput;
      if (!this.showCommentInput) {
        localStorage.removeItem((this.newText))
        this.newText = ""
      }
    }
  }
}
</script>

<template>

  <div class="d-flex position-relative">
    <div class="d-flex position-absolute top-0 end-0 mt-3">
      <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
    </div>
  </div>
  <div>
    <div style="display: flex; justify-content: center;">
      <h2 class="h2">Profile</h2>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>


  </div>

  <div class="col-lg-7">
    <div class="d-flex align-items-center">
      <h1 style="white-space: nowrap;">{{ userProfile.profile.username }}</h1>
      <div class="col-lg-6 mx-5">
        <div class="d-flex mt-3 justify-content-between">
          <h6 style="margin-right: 10px;">{{ userProfile.profile.follower_count}} Follower</h6>
          <h6 style="margin-right: 10px;">{{ userProfile.profile.following_count }} Following</h6>
          <h6>{{ userProfile.profile.post_count }} Post</h6>
        </div>
      </div>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
  </div>

  <div class="post-grid">
    <div v-for="post in userProfile.posts" :key="post.ID" class="post-container">
      <img :src="'data:image/jpeg;base64,' + post.file" alt="Post Image" class="post-image img-fluid">
      <p><span style="font-weight: bold;">{{ username }}</span>: {{ post.description }}</p>
      <button type="button" class="btn" @click="showLikes(post.ID)">
        Likes: {{ post.like_count }}
      </button>
      <button type="button" class="btn" @click="toggleLike(post.ID)">
        <svg class="heart-icon" viewBox="0 0 35 35" xmlns="http://www.w3.org/2000/svg">
          <path d="M7.975 2c-2.235.116-4.365 1.203-5.82 2.89C.7 6.57 0 8.786 0 11c0 1.938.697 3.816 1.646 5.46.95 1.644 2.19 3.077 3.5 4.394 2.824 2.833 6.08 5.232 9.622 7.09.145.076.32.076.464 0 3.543-1.858 6.798-4.257 9.622-7.09 1.31-1.317 2.55-2.75 3.5-4.393C29.304 14.817 30 12.94 30 11c0-2.22-.7-4.428-2.154-6.11C26.39 3.202 24.26 2.115 22.026 2c-1.516-.078-3.045.286-4.362 1.04-1.097.626-1.975 1.558-2.664 2.614-.69-1.056-1.567-1.988-2.664-2.615C11.02 2.285 9.49 1.92 7.976 2zm.05 1c1.32-.068 2.665.25 3.813.906 1.148.656 2.107 1.652 2.72 2.824.186.36.698.36.885 0 .612-1.172 1.57-2.168 2.72-2.824 1.147-.656 2.49-.974 3.812-.906 1.942.1 3.837 1.062 5.115 2.54C28.37 7.023 29 9 29 11c0 1.73-.628 3.43-1.512 4.96-.885 1.535-2.064 2.904-3.342 4.186-2.686 2.697-5.788 4.975-9.146 6.766-3.358-1.79-6.46-4.07-9.146-6.766-1.278-1.282-2.457-2.65-3.342-4.185C1.628 14.43 1 12.73 1 11c0-2 .63-3.978 1.91-5.46C4.188 4.063 6.083 3.1 8.025 3z"/>
        </svg>
      </button>
      <div>
        <button type="button" class="btn" @click="showComments(post.ID)">
          Comments: {{ post.comment_count }}
        </button>
        <button type="button" class="btn" @click="toggleCommentInput(post.ID)">
          <svg class="keyboard-icon" viewBox="0 0 35 35" xmlns="http://www.w3.org/2000/svg">
            <path d="M6.21,13.29a.93.93,0,0,0-.33-.21,1,1,0,0,0-.76,0,.9.9,0,0,0-.54.54,1,1,0,1,0,1.84,0A1,1,0,0,0,6.21,13.29ZM13.5,11h1a1,1,0,0,0,0-2h-1a1,1,0,0,0,0,2Zm-4,0h1a1,1,0,0,0,0-2h-1a1,1,0,0,0,0,2Zm-3-2h-1a1,1,0,0,0,0,2h1a1,1,0,0,0,0-2ZM20,5H4A3,3,0,0,0,1,8v8a3,3,0,0,0,3,3H20a3,3,0,0,0,3-3V8A3,3,0,0,0,20,5Zm1,11a1,1,0,0,1-1,1H4a1,1,0,0,1-1-1V8A1,1,0,0,1,4,7H20a1,1,0,0,1,1,1Zm-6-3H9a1,1,0,0,0,0,2h6a1,1,0,0,0,0-2Zm3.5-4h-1a1,1,0,0,0,0,2h1a1,1,0,0,0,0-2Zm.71,4.29a1,1,0,0,0-.33-.21,1,1,0,0,0-.76,0,.93.93,0,0,0-.33.21,1,1,0,0,0-.21.33A1,1,0,1,0,19.5,14a.84.84,0,0,0-.08-.38A1,1,0,0,0,19.21,13.29Z"/>
          </svg>
        </button>
        <div v-if="showCommentInput" style="margin-right: 10px;">
          <input type="text" id="newDescription" v-model="newText" class="form-control form-control-sm" placeholder="comment text" aria-label="Recipient's comment" aria-describedby="basic-addon2">
          <button v-if="toggleCommentInput" type="button" class="btn btn-sm btn-primary" @click="commentPhoto(post.ID)">Send</button>
        </div>
      </div>
      <div class="delete-button-container">
        <button type="button" class="btn delete-button" @click="deletePhoto(post.ID)">
          X
        </button>
      </div>
    </div>
  </div>

  <div class="liked-users-overlay" v-if="this.showLikeWindow">
    <div class="liked-users-modal">
      <h2>Likes</h2>
      <ul>
        <li v-for="user in this.fullPost.like_owners" :key="user.username">{{ user.username }}</li>
      </ul>
      <button @click="this.closeLikeWindow()">Close</button>
    </div>
  </div>

  <div class="liked-users-overlay" v-if="this.showCommentWindow">
    <div class="liked-users-modal">
      <h2>Comments</h2>
      <ul>
        <li v-for="comment in this.fullPost.comments" :key="comment.owner_ID">{{ comment.owner_ID + ": " + comment.text }}</li>
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
              <RouterLink :to="'/users/' + userProfile.profile.ID + '/profile'" class="nav-link" style="font-size: 20px;" @click="getMyProfile">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
                Profile
              </RouterLink>
            </li>
            <li class="nav-item">
              <button type="button" class="btn btn-sm" style="font-size: 12px;" @click="togglePhotoInput">
                Upload
              </button>
              <div v-if="showUploadInput" style="margin-right: 10px;">
                <input type="file" id="newPhoto" @change="handleFileChange" class="form-control form-control-sm">
                <input type="text" id="newDescription" v-model="newText" class="form-control form-control-sm" placeholder="Photo description" aria-label="Recipient's description" aria-describedby="basic-addon2">
              </div>
              <button v-if="showUploadInput" type="button" class="btn btn-sm btn-primary" @click="uploadPhoto">Upload</button>
            </li>
            <li class="nav-item">
              <button type="button" class="btn btn-sm" style="font-size: 12px;" @click="toggleUsernameInput">
                Change Username
              </button>
              <div class="d-flex ">
                <input v-if="showUsernameInput" type="text" id="newUsername" v-model="newUsername" class="form-control form-control-sm"
                       placeholder="New Username" aria-label="Recipient's username" aria-describedby="basic-addon2">
                <button v-if="showUsernameInput" type="button" class="btn btn-sm btn-primary ml-2 me-2" @click="setMyUserName">
                  Change
                </button>
              </div>
            </li>
          </ul>
          <div class="mt-auto mb-3">
            <RouterLink to="/" class="nav-link" style="font-size: 20px;" @click="doLogout">
              <svg fill="#000000" width="20px" height="20px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="M2.293,11.293l4-4A1,1,0,1,1,7.707,8.707L5.414,11H17a1,1,0,0,1,0,2H5.414l2.293,2.293a1,1,0,1,1-1.414,1.414l-4-4a1,1,0,0,1,0-1.414ZM20,4V20a1,1,0,0,0,2,0V4a1,1,0,0,0-2,0Z"/></svg>
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
  display: grid ;
  grid-template-columns: repeat(auto-fill, minmax(calc(33.33% - 15px), 1fr));
  grid-gap: 15px;
}

.post-container {
  position: relative;
  border: 2px solid #ccc;
  margin-bottom: 20px;
  padding: 10px;
  width: 100%;
  margin-top: 50px;
}

.post-image {
  margin-bottom: 10px;
}

.delete-button-container {
  position: absolute;
  bottom: 0;
  right: 0;
}

.delete-button {
  font-size: 20px;
  color: red;
}

.heart-icon {
  width: 15px;
  height: 15px;
  fill: red;
}
.keyboard-icon {
  width: 20px;
  height: 20px;
  fill: black;
}

.liked-users-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.liked-users-modal {
  background-color: white;
  padding: 20px;
  border-color: #000000;
  border-width: 2px;
  border-style: solid;
  max-height: 80%;
  overflow-y: auto;
}
</style>
