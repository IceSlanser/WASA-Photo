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
        followers: []
      },

      stream: [
        {
          post: {
            ID: 0,
            profile_username: "",
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
      showSearchInput: false,
      showCommentWindow: false,
      showLikeWindow: false,
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
      try {
        let response = await this.$axios.get(`/stream`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.stream = response.data
        if (this.stream == null) {
          this.stream = []
        }
        await localStorage.setItem("stream", JSON.stringify(this.stream));
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to get stream.";
        } else if (e.response && e.response.status === 404) {
          this.error = "Stream not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later.";
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

  <div class="post-grid">
    <div v-for="(post, index) in sortedPosts" :key="post.ID" class="post-container">
      <img v-if="post.post.file" :src="'data:image/jpeg;base64,' + post.post.file" alt="Post Image" class="post-image img-fluid align-content-center">
      <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
      <div class="d-flex justify-content-between">
        <p><span style="font-weight: bold;">{{ post.post.profile_ID }}</span>: {{ post.post.description }}</p>
        <p>{{post.post.date_time}}</p>
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
      <div class="delete-button-container" v-if="this.myUsername === this.userProfile.profile.username">
        <button type="button" class="btn delete-button" @click="deletePhoto(post.ID)">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
        </button>
      </div>
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
</style>
