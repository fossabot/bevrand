﻿using bevrand.authenticationapi.DAL.Models;
using bevrand.authenticationapi.Services;
using Microsoft.AspNetCore.JsonPatch;
using Microsoft.AspNetCore.Mvc;

namespace bevrand.authenticationapi.Controllers
{
    [Route("api/[controller]")]
    public class UsersController : Controller
    {

       private readonly IUsersLogic _usersLogic;

       public UsersController(IUsersLogic usersLogic)
       {
           _usersLogic = usersLogic;
       }


        [HttpGet]
        public IActionResult Get()
        {
            var result = _usersLogic.GetAllUsersFromDataBase();
            return Ok(result);
        }

        [HttpGet("{id}")]
        public IActionResult GetById(int id)
        {
            var result = _usersLogic.GetById(id);
            return Ok(result);
        }

        [HttpGet("by-email/{emailaddress}", Name = "GetByEmail")]
        public IActionResult GetByEmail(string emailaddress)
        {
            var result = _usersLogic.GetByEmailAddress(emailaddress);
            return Ok(result);
        }
        
        
        [HttpGet("by-username/{username}", Name = "GetByUserName")]
        public IActionResult GetByUserName(string username)
        {
            var result = _usersLogic.GetByUserName(username);
            return Ok(result);
        }
        
        [HttpPost]
        public IActionResult Create([FromBody] PostUserModel user)
        {
           _usersLogic.CreateANewUser(user);
            var userToReturn = _usersLogic.GetByUserName(user.UserName);
            return CreatedAtRoute("GetByUserName", new {username = user.UserName}, userToReturn);
        }

        [HttpPatch]
        public IActionResult Patch([FromQuery]int id, [FromBody]JsonPatchDocument<PatchUserModel> user)
        {
            var result = _usersLogic.PatchAUser(id, user);
            return Ok(result);

        }

        [HttpDelete]
        public IActionResult Delete([FromQuery]int id)
        {
            _usersLogic.DeleteAUser(id);
            return new NoContentResult();
        }
    }
}



