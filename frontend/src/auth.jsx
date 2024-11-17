import React, { useState } from "react";
import axios from "axios";

// You will need to replace this with actual routes based on your backend API setup.
const API_BASE_URL = "http://localhost:8080"; // Your backend base URL

const Login = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const handleEmailLogin = async () => {
        try {
            const response = await axios.post(`${API_BASE_URL}/login/email`, { email, password });
            console.log("Login Success:", response.data);
        } catch (error) {
            console.error("Login Failed:", error);
        }
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
            {/* Login UI */}
            <div className="bg-white p-8 rounded shadow-md w-full max-w-md">
                <h2 className="text-2xl font-bold mb-6 text-center">Login</h2>
                <input
                    type="email"
                    placeholder="Enter Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full p-2 mb-4 border border-gray-300 rounded"
                />
                <input
                    type="password"
                    placeholder="Enter Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full p-2 mb-4 border border-gray-300 rounded"
                />
                <button onClick={handleEmailLogin} className="w-full py-2 mb-4 bg-blue-500 text-white rounded">
                    Login with Email
                </button>
                <div className="text-center mb-4">or</div>
                <button className="w-full py-2 mb-2 bg-yellow-500 text-white rounded flex items-center justify-center">
                    <i className="fab fa-ethereum mr-2"></i> Connect with MetaMask
                </button>
                <button className="w-full py-2 bg-purple-500 text-white rounded flex items-center justify-center">
                    <i className="fab fa-ethereum mr-2"></i> Connect with Solana
                </button>
                <div className="text-center mt-4">
                    <a href="#" className="text-blue-500">Don't have an account? Sign Up</a>
                </div>
            </div>
        </div>
    );
};

const SignUp = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const handleEmailSignUp = async () => {
        try {
            const response = await axios.post(`${API_BASE_URL}/signup/email`, { email, password });
            console.log("Sign Up Success:", response.data);
        } catch (error) {
            console.error("Sign Up Failed:", error);
        }
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
            {/* SignUp UI */}
            <div className="bg-white p-8 rounded shadow-md w-full max-w-md">
                <h2 className="text-2xl font-bold mb-6 text-center">Sign Up</h2>
                <input
                    type="email"
                    placeholder="Enter Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full p-2 mb-4 border border-gray-300 rounded"
                />
                <input
                    type="password"
                    placeholder="Enter Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full p-2 mb-4 border border-gray-300 rounded"
                />
                <button onClick={handleEmailSignUp} className="w-full py-2 mb-4 bg-blue-500 text-white rounded">
                    Sign Up with Email
                </button>
                <div className="text-center mb-4">or</div>
                <button className="w-full py-2 mb-2 bg-yellow-500 text-white rounded flex items-center justify-center">
                    <i className="fab fa-ethereum mr-2"></i> Connect with MetaMask
                </button>
                <button className="w-full py-2 bg-purple-500 text-white rounded flex items-center justify-center">
                    <i className="fab fa-ethereum mr-2"></i> Connect with Solana
                </button>
                <div className="text-center mt-4">
                    <a href="#" className="text-blue-500">Already have an account? Login</a>
                </div>
            </div>
        </div>
    );
};

const Auth = () => {
    const [isLogin, setIsLogin] = useState(true);

    return (
        <div>
            {isLogin ? <Login /> : <SignUp />}
            <div className="fixed bottom-4 right-4">
                <button 
                    onClick={() => setIsLogin(!isLogin)} 
                    className="px-4 py-2 bg-blue-500 text-white rounded">
                    {isLogin ? 'Go to Sign Up' : 'Go to Login'}
                </button>
            </div>
        </div>
    );
};

export default Auth;
