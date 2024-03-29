import { Component, Inject, OnDestroy } from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import { Subscription } from 'rxjs';
import { UserService } from '../user.service';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import { User } from '../user';
import { InstantErrorStateMatcher } from './error-state.matcher';
import {MatSelectModule} from '@angular/material/select';

@Component({
  selector: 'user-dialog',
  templateUrl: './user-dialog.component.html',
  styleUrls: ['./user-dialog.component.scss']
})
export class UserDialog implements OnDestroy {
  controlGroup: FormGroup;
  errorStateMatcher = new InstantErrorStateMatcher();
  addSubscription!: Subscription;
  updateSubscription!: Subscription;
  deleteSubscription!: Subscription;
  statuses: string[] = ['Active', 'Inactive', 'Terminated']

  constructor(
    @Inject(MAT_DIALOG_DATA) public user: User,
    public dialogRef: MatDialogRef<UserDialog>,
    public service: UserService
  ) {
    this.controlGroup = new FormGroup({
      user_name: new FormControl(user.user_name, Validators.required),
      first_name: new FormControl(user.first_name, Validators.required),
      last_name: new FormControl(user.last_name, Validators.required),
      email: new FormControl(user.email, Validators.required),
      status: new FormControl(user.user_status, Validators.required),
      department: new FormControl(user.department, Validators.required),
    });
  }

  save(): void {
    this.user.user_name = this.formValue('username');
    this.user.first_name = this.formValue('first_name');
    this.user.last_name = this.formValue('last_name');
    this.user.email = this.formValue('email');
    this.user.user_status = this.formValue('status');
    this.user.department = this.formValue('department')

    if (!this.user.user_id) {
      this.addSubscription = this.service.add(this.user)
        .subscribe(this.dialogRef.close);
    } else {
      this.updateSubscription = this.service.update(this.user)
        .subscribe(this.dialogRef.close);
    }
  }

  delete(): void {
    this.deleteSubscription = this.service.delete(this.user.user_id!)
      .subscribe(this.dialogRef.close);
  }

  hasError(controlName: string, errorCode: string): boolean {
    return !this.controlGroup.valid && this.controlGroup.hasError(errorCode, [controlName]);
  }

  ngOnDestroy(): void {
    if (this.addSubscription) {
      this.addSubscription.unsubscribe();
    }
    if (this.updateSubscription) {
      this.updateSubscription.unsubscribe();
    }
    if (this.deleteSubscription) {
      this.deleteSubscription.unsubscribe();
    }
  }

  private formValue(controlName: string): any {
      return this.controlGroup.get(controlName)?.value
  }
}


