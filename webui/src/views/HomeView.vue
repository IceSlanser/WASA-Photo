<script>
import {RouterLink} from "vue-router";

export default {
  components: {RouterLink},
  data: function() {
    return {
      error: null,
      userProfile: {
        profile: {
          ID: localStorage.getItem("ID"),
          username: localStorage.getItem("username"),
          following_count: 0,
          follower_count: 0,
          post_count: 0
        },
        posts: [
          {
            ID: 0,
            profile_ID: 0,
            file: [],
            description: "",
            like_count: 0,
            comment_count: 0,
            date_time: Date
          }
        ],
        followings: [
          {
            ID: 0
          }
        ],
        followers: [
          {
            ID: 0
          }
        ]
      }
    }
  },
	methods: {
    async doLogout() {
      localStorage.clear()
      this.$router.push({ path : '/'})
    },
    async uploadPhoto() {

    },
    async getMyProfile() {
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
          this.err = "Failed to request user's profile.";
        } else if (e.response && e.response.status === 404) {
          this.err = "User not found.";
        } else if (e.response && e.response.status === 500) {
          this.err = "An internal error occurred, please try again later.";
        } else {
          this.err = e.toString();
        }

      }
    }
	},
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
</style>
