import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { StakeholdersService } from '../../../core/services/stakeholders.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './profile.component.html'
})
export class ProfileComponent implements OnInit {
  user: any = {}; // moze se definisati interface za tipizaciju
  editData: any = {};
  isEditing = false;

  constructor(private stakeholdersService: StakeholdersService) {}

  ngOnInit(): void {
    console.log('Token:', localStorage.getItem('access_token'));

    this.stakeholdersService.getProfile().subscribe({
      next: (data) => {
        this.user = data;
        console.log('User data:', this.user); // Debug ceo objekat
        console.log('Profile image:', this.user.profileImage); // Debug specifično sliku
      },
      error: (err) => {
        console.error('Failed to load profile', err);
      }
    });
  }

  startEdit(): void {
    this.isEditing = true;
    this.editData = {
      name: this.user.name || '',
      surname: this.user.surname || '',
      profile_image: this.user.profile_image || '',
      biography: this.user.biography || '',
      motto: this.user.motto || ''
    };
  }

  cancelEdit(): void {
    this.isEditing = false;
    this.editData = {};
  }

  saveProfile(): void {
    this.stakeholdersService.updateProfile(this.editData).subscribe({
      next: (data) => {
        this.user.name = this.editData.name;
        this.user.surname = this.editData.surname;
        this.user.profile_image = this.editData.profile_image;
        this.user.biography = this.editData.biography;
        this.user.motto = this.editData.motto;

        this.isEditing = false;
        console.log('Profile updated successfully');
        alert('Profile updated successfully!');
      },
      error: (err) => {
        console.error('Failed to update profile', err);
        alert('Failed to update profile. Please try again.');
      }
    });
  }

}