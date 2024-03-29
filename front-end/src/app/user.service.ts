import { Injectable } from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {Observable} from "rxjs";
import {User} from "./user";

@Injectable({
  providedIn: 'root'
})
export class UserService {
  baseUrl: string = 'http://localhost:8080/api/v1/users';
  readonly headers = new HttpHeaders()
    .set('Content-Type', 'application/json');

  constructor(private http: HttpClient) { }

  getAll(): Observable<User[]> {
    return this.http.get<User[]>(this.baseUrl, {headers: this.headers});
  }

  add(user: User): Observable<User> {
    return this.http.post<User>(this.baseUrl, user, {headers: this.headers});
  }

  update(user: User): Observable<User> {
    return this.http.put<User>(
      `${this.baseUrl}/${user.user_id}`, user, {headers: this.headers}
    );
  }

  delete(id: number): Observable<User> {
    return this.http.delete<User>(`${this.baseUrl}/${id}`, {headers: this.headers});
  }
}
