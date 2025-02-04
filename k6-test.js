import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: "5s", target: 1000 },
        { duration: "20s", target: 1000 },
        { duration: "5s", target: 0 },
    ],
};

export default () => {
    let Data_Login = JSON.stringify({
        email: "akromjonotaboyev@gmail.com",
        password: "Akrom2005",
        platform: "admin",
    });

    // let uniqueId = Math.random().toString(36).substring(2, 8);
    // let emailDomain = "gmail.com";

    // let Data_Register = JSON.stringify({
    //     full_name: "Test",
    //     user_type: "user",
    //     user_role: "user",
    //     username: "testusername",
    //     email: `test+${uniqueId}@${emailDomain}`,
    //     profile_picture: `${uniqueId}`,
    //     status: "inverify",
    //     password: "1234",
    //     gender: "male",
    // });
    

    // let loginParams = {
    //     headers: {
    //         "Content-Type": "application/json"
    //     }
    // };

    

    let registerParams = {
        headers: {
            "Authorization": `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwbGF0Zm9ybSI6ImFkbWluIiwic2Vzc2lvbl9pZCI6IjMzMjdjZThiLWJhN2UtNGEwZi04MWM5LWVkZWQwOTY2MTQ1NCIsInN1YiI6IjcyY2JlM2ZiLTAyY2QtNGZhMy04YWM0LTU0NDMyMjhmOTc4ZiIsInVzZXJfcm9sZSI6InN1cGVyYWRtaW4iLCJ1c2VyX3R5cGUiOiJhZG1pbiJ9.J9yJXKkYlcHfuFDFEhvDXgKAdBIZlahmfEWgRkgqOXM`,
            "Content-Type": "application/json",
        }
    };


  

    // const resRegister = http.post('http://localhost:8080/v1/user/', Data_Register, registerParams);

    // check(resRegister, {
    //     "status code 201": (r) => r.status === 201
    // });

    const resGetSingleUser = http.get(`http://localhost:9090/v1/user/72cbe3fb-02cd-4fa3-8ac4-5443228f978f`, registerParams);

    check(resGetSingleUser, {
        "status code 200": (r) => r.status === 200
    });

    // console.log('Register Response:', resRegister.body);
    // console.log('Get Single User Response:', resGetSingleUser.body);
};

