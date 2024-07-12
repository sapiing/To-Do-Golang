"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";

export default function DashboardPage() {
  const [tasks, setTasks] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const router = useRouter();
  const [isLoggedOut, setIsLoggedOut] = useState(false);

  const handleLogout = async () => {
    try {
      const response = await fetch("http://localhost:8080/logout", {
        method: "POST",
        credentials: "include", // Important for sending cookies
      });

      if (response.ok) {
        setIsLoggedOut(true);
        localStorage.removeItem("token"); // Remove token from local storage
        router.push("/login"); // Redirect to login page (or another suitable page)
      } else {
        console.error("Logout failed:", response.statusText);
      }
    } catch (error) {
      console.error("Error during logout:", error);
    }
  };

  const fetchTasks = useCallback(async () => {
    try {
      const token = localStorage.getItem("token"); // Get token from local storage
      const response = await fetch("http://localhost:8080/api/task", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
      }

      const data = await response.json();

      const formattedTasks = data.map((task) => ({
        ...task,
        date: new Date(task.date).toISOString(), // Adjust to your preferred date format
      }));
      setTasks(formattedTasks);

      setTasks(data);
    } catch (error) {
      setError(error.message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  async function handleToggleComplete(task) {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(
        `http://localhost:8080/api/task/${task.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            id: task.id,
            completed: !task.completed,
          }), // Toggle completed
        }
      );

      if (!response.ok) {
        throw new Error("Failed to update task status");
      }

      // Update the task list after successful status change
      fetchTasks();
    } catch (error) {
      setError(error.message);
    }
  }

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]); // Empty dependency array ensures this runs only once on mount

  const goToDetailPage = () => {
    router.push("/detail");
  };

  return (
    <div>
      <div className="d-flex justify-content-space-between">
        <button
          onClick={goToDetailPage}
          className="bg-teal-500 hover:bg-teal-600 transition-colors duration-300 text-white font-semibold py-2.5 px-6 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-teal-400"
        >
          Task Detail
        </button>
        <button
          onClick={handleLogout}
          className="bg-teal-500 hover:bg-teal-600 transition-colors duration-300 text-white font-semibold py-2.5 px-6 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-teal-400"
        >
          Logout
        </button>
      </div>

      {isLoading ? (
        <p className="text-gray-500 mt-4">Loading tasks...</p>
      ) : error ? (
        <p className="text-red-500 mt-4">Error: {error}</p>
      ) : (
        <table className="table-auto w-full mt-6 border-collapse">
          <thead className="text-left bg-gray-100">
            <tr className="rounded-lg">
              <th className="py-3 px-4 font-semibold rounded-tl-lg">#</th>
              <th className="py-3 px-4 font-semibold rounded-tl-lg">Title</th>
              <th className="py-3 px-4 font-semibold rounded-tl-lg">Description</th>
            </tr>
          </thead>
          <tbody>
            {tasks.length === 0 ? (
              <tr>
                <td colSpan="4" className="py-4 text-center text-gray-500">
                  No tasks found. Create a new task to get started!
                </td>
              </tr>
            ) : (
              tasks.map((task, index = 1) => (
                <tr key={task.id} className="hover:bg-gray-50">
                  <td className="py-2 px-4 border-b">{index + 1}</td>
                  <td className="py-2 px-4 border-b">{task.title}</td>
                  <td className="py-2 px-4 border-b">{task.description}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      )}
    </div>
  );
}
