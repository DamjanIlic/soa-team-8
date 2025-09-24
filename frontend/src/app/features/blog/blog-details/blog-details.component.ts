import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { MarkdownModule } from 'ngx-markdown';
import { BlogService } from '../../../core/services/blog.service';
import { Blog } from '../../../core/models/blog.model';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms'; // za ngModel
import { Comment } from '../../../core/models/comment.model';
import { FollowService } from '../../../core/services/follow.service';
import { StakeholdersService } from '../../../core/services/stakeholders.service';

@Component({
  selector: 'app-blog-details',
  standalone: true,
  imports: [
    CommonModule,
    HttpClientModule,
    MarkdownModule,
    FormsModule
  ],
  templateUrl: './blog-details.component.html',
  styleUrls: ['./blog-details.component.css']
})
export class BlogDetailsComponent implements OnInit {
  blog?: Blog;
  blogId!: string;
  blogOwnerId?: string;

  comments: Comment[] = [];
  user_id = "";
  isMyBlog = false;
  newCommentText = '';
  followingMap: Record<string, boolean> = {};
  a = 5
  currentUser: any;
  // tmp jwt
  constructor(
    private route: ActivatedRoute,
    private blogService: BlogService,
    private followService: FollowService,
    private stakeholdersService: StakeholdersService
  ) {}

  ngOnInit(): void {
    this.blogId = this.route.snapshot.paramMap.get('id')!;
    this.loadCurrentUser();
    const userId = this.getCurrentUserId();
    if (userId){
      this.user_id = userId
      console.log('ui', this.user_id)
    }
    // this.currentUser = { user_id: userId };
    // console.log(this.currentUser)
    this.loadBlog();
    this.loadComments();
    this.loadFollowing();
  }
  loadCurrentUser(): void {
    this.stakeholdersService.getProfile().subscribe({
      next: (user) => {
        this.currentUser = user; 
        console.log('Current user:', this.currentUser);
        console.log(this.user_id)
      },
      error: (err) => {
        console.error('Failed to load current user', err);
      }
    });
  }
  private getCurrentUserId(): string | null {
    const token = localStorage.getItem('access_token'); // gde čuvaš JWT
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.user_id; // ili kako se polje zove u JWT
    } catch (e) {
      console.error('Invalid token', e);
      return null;
    }
  }

  loadBlog(): void {
    this.blogService.getById(this.blogId).subscribe({
      next: (data) => {
        this.blog = data
        console.log(this.blog)
        this.blogOwnerId = data.user_id;
        console.log('bi', this.blogOwnerId)
        this.isMyBlog = this.blogOwnerId === this.user_id;
      },
      error: (err) => console.error(err)
    });
  }

  loadComments(): void {
    this.blogService.getComments(this.blogId).subscribe({
      next: (data) => this.comments = data,
      error: (err) => console.error(err)
    });
  }

  likeBlog(): void {
    this.blogService.like(this.blogId).subscribe({
      next: (likes) => { if (this.blog) this.blog.likes = likes; },
      error: (err) => console.error(err)
    });
  }

  unlikeBlog(): void {
    this.blogService.unlike(this.blogId).subscribe({
      next: (likes) => { if (this.blog) this.blog.likes = likes; },
      error: (err) => console.error(err)
    });
  }

  addComment(): void {
    if (!this.newCommentText.trim()) return;

    const comment: Comment = {
      username: 'You',
      blogId: this.blogId,
      text: this.newCommentText.trim()
    };

    this.blogService.createComment(this.blogId, comment).subscribe({
      next: (c) => {
        c.username = 'You';
        this.comments.push(c);
        this.newCommentText = '';
      },
      error: (err) => console.error(err)
    });
  }

  
  followUser(targetUserId: string) {
    if (!targetUserId) {
      console.warn('No userId to follow');
      return;
    }

    this.followService.followUser(targetUserId).subscribe({
      next: (res) => {
        console.log('Successfully followed user', res);
        // Update local map
        this.followingMap[targetUserId] = true;
      },
      error: (err) => {
        console.error('Failed to follow user', err);
      }
    });
  }

  unfollowUser(targetUserId: string) {
    if (!targetUserId) {
      console.warn('No userId to unfollow');
      return;
    }

    this.followService.unfollowUser(targetUserId).subscribe({
      next: (res) => {
        console.log('Successfully unfollowed user', res);
        // Update local map
        delete this.followingMap[targetUserId];
      },
      error: (err) => {
        console.error('Failed to unfollow user', err);
      }
    });
  }
  loadFollowing(): void {
    this.followService.getFollowing().subscribe({
      next: (res: { following: string[] }) => {
        this.followingMap = {};
        res.following.forEach(id => this.followingMap[id] = true);
      },
      error: (err) => console.error('Failed to load following', err)
    });
  }

  isFollowing(userId: string): boolean {
    if (userId === 'test') {
      return false;
    }
    if (userId === 'testF'){
      return true;
    }
    console.log(this.followingMap)
    return this.followingMap[userId];
  }
}
