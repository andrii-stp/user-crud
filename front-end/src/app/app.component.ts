import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from "@angular/material/dialog";
import { UserDialog } from "./user-dialog/user-dialog.component";
import { User } from './user';
import { Subscription } from 'rxjs';
import { UserService } from './user.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {
  displayedColumns = ['Username', 'FirstName', 'LastName', 'Email', 'Status', 'Department'];
  dataSource: User[] = [];
  getAllSubscription!: Subscription;
  dialogSubscription!: Subscription;

  constructor(public dialog: MatDialog, public service: UserService) {}

  openEditDialog(user: User) {
    this.openDialog(new User(user.user_id, user.user_name, user.first_name, user.last_name, user.email, user.user_status, user.department));
  }

  openNewDialog(): void {
    this.openDialog(new User());
  }

  private openDialog(user: User): void {
    this.dialogSubscription = this.dialog
      .open(UserDialog, {data: user, minWidth: '30%'})
      .afterClosed().subscribe(() => this.loadUserList());
  }

  private loadUserList(): void {
    this.getAllSubscription = this.service.getAll()
      .subscribe(users => this.dataSource = users);
  }

  ngOnInit(): void {
    this.loadUserList();
  }

  ngOnDestroy(): void {
    this.getAllSubscription.unsubscribe();
    if (this.dialogSubscription) {
      this.dialogSubscription.unsubscribe();
    }
  }
}
