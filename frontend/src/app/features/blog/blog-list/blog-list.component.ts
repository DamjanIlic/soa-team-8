import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BlogService } from '../../../core/services/blog.service'; 
import { Blog } from '../../../core/models/blog.model'; 
import { catchError, of } from 'rxjs';
import { RouterModule } from '@angular/router';
import { FollowService } from '../../../core/services/follow.service';


interface RecommendedUser {
  user_id: string;
  username: string;
  image: string;
  motto: string;
}

interface RecommendationsResponse {
  recommendations: RecommendedUser[];
}

@Component({
  selector: 'app-blog-list',
  standalone: true,
  imports: [CommonModule, RouterModule], // HttpClientModule više nije potrebno
  templateUrl: './blog-list.component.html',
  styleUrls: ['./blog-list.component.css']
})


export class BlogListComponent implements OnInit {
  blogs: Blog[] = [];
  blogsForUser: Blog[] = [];
  loading = false;
  error = '';
  selectedTab: 'foryou' | 'explore' | 'recommend'= 'foryou';
  showRecommendForm = false;
  recommendations: any[] = [];
  loadingRecommendations = false;
  constructor(private blogService: BlogService, private followService: FollowService) {}

  ngOnInit(): void {
    this.loadForYouBlogs();
    //this.fetchBlogs();
  }

  fetchBlogs() {
    this.loading = true;
    this.blogService.getAll()
      .pipe(
        catchError(err => {
          this.error = 'Error fetching blogs.';
          console.error(err);
          return of([]);
        })
      )
      .subscribe(data => {
        this.blogs = data;
        this.loading = false;
      });
  }

  loadForYouBlogs() {
    this.selectedTab = 'foryou';
    this.loading = true;
    this.blogService.getForCurrentUser() // nova funkcija u blogService
      .pipe(
        catchError(err => {
          this.error = 'Error fetching For You blogs.';
          console.error(err);
          return of([]);
        })
      )
      .subscribe(data => {
        this.blogsForUser = data;
        this.blogs = this.blogsForUser;
        this.loading = false;
      });
  }

  loadExploreBlogs() {
    this.selectedTab = 'explore';
    this.fetchBlogs(); // koristi postojeću getAll funkciju
  }
  openForm() {
    console.log('h')
    this.showRecommendForm = true;
    this.loadRecommendations();
  }

  loadRecommendations() {
    console.log('load r')
    this.loadingRecommendations = true;
    this.followService.getRecommendations().subscribe({
      next: (res: any) => {
        this.recommendations = res.recommendations || [];
        console.log(this.recommendations)
        this.loadingRecommendations = false;
      },
      error: (err) => {
        console.error('Greška prilikom učitavanja preporuka', err);
        this.loadingRecommendations = false;
      }
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

        // Ukloni korisnika iz liste preporuka
        this.recommendations = this.recommendations.filter(
          user => user.user_id !== targetUserId
        );
      },
      error: (err) => {
        console.error('Failed to follow user', err);
      }
    });
  }
}
