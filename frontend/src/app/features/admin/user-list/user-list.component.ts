import { Component } from '@angular/core';
import { NgFor, NgClass, CommonModule } from '@angular/common';
import {  OnInit } from '@angular/core';
import { StakeholdersService } from '../../../core/services/stakeholders.service'; 

interface User {
  id: string;         // Stakeholder ID
  userId: string;     // User ID iz users tabele
  name: string;
  surname: string;
  username: string;
  email: string;
  role: string;
  status: string;
  profile_image?: string;
  blocked: boolean;
}

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [NgFor, NgClass, CommonModule],
  templateUrl: './user-list.component.html',
})
export class UserListComponent {
  users: User[] = [];
  loading = true;
  error: string | null = null;

  constructor(private stakeholdersService: StakeholdersService) {}

  ngOnInit(): void {
    this.loadUsers();
  }

loadUsers(): void {
  this.stakeholdersService.getAllUsers().subscribe({
    next: (data: any[]) => {
      this.users = data.map(u => ({
        id: u.id,
        userId: u.user_id,
        name: u.name,
        surname: u.surname,
        username: u.username,
        email: u.email,
        role: u.role,
        blocked: u.blocked,
        avatarUrl: u.profile_image || 'https://via.placeholder.com/64',
        status: u.blocked ? 'Suspended' : 'Active', // OBAVEZNO dodaj status
      }));
      this.loading = false;
    },
    error: (err) => {
      this.error = 'Failed to load users';
      console.error(err);
      this.loading = false;
    }
  });
}


  blockUser(userId: string): void {
    this.stakeholdersService.blockUser(userId).subscribe({
      next: () => this.loadUsers(),
      error: (err) => console.error(err)
    });
  }

  unblockUser(userId: string): void {
    this.stakeholdersService.unblockUser(userId).subscribe({
      next: () => this.loadUsers(),
      error: (err) => console.error(err)
    });
  }
}
