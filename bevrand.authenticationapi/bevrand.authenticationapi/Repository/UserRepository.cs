﻿using System.Collections.Generic;
using System.Linq;
using bevrand.authenticationapi.Data;
using bevrand.authenticationapi.DAL.Models;
using bevrand.authenticationapi.Repository.Models;
using Microsoft.EntityFrameworkCore;

namespace bevrand.authenticationapi.Repository
{


    public class UserRepository : IUserRepository
    {
        private UserContext _context;
        
        public UserRepository(UserContext context)
        {
            _context = context;
        }

        public UserModel Add(UserModel user)
        {
            _context.Users.Add(user);
            _context.SaveChanges();
            return user;
        }

        public void Update(UserModel user)
        {
            _context.Users.Update(user);
            _context.SaveChanges();
        }

        public void Delete(UserModel user)
        {
            _context.Users.Remove(user);
            _context.SaveChanges();
        }
        
        public bool CheckIfUserExists(string name)
        {
            return _context.Users.Any(u => u.UserName == name);
        }
        
        public bool CheckIfEmailExists(string name)
        {
            return _context.Users.Any(u => u.EmailAddress == name);
        }
        
        public bool CheckIfIdExists(int id)
        {
            return _context.Users.Any(u => u.Id == id);
        }

        public IEnumerable<UserModel> GetAllUsers()
        {
            return _context.Users.OrderBy(u => u.Id);
        }

        public UserModel GetSingleUser(int id)
        {
            return _context.Users.AsNoTracking().FirstOrDefault(u => u.Id == id);
        }
        
        public UserModel GetSingleUser(string userName)
        {
            return _context.Users.FirstOrDefault(u => u.UserName == userName);
        }

        public UserModel GetSingleUserEmail(string emailAddress)
        {
            return _context.Users.FirstOrDefault(u => u.EmailAddress == emailAddress);
        }


    }
}


