function ForgotPass() {
  return (
    <div>
      <div className="flex justify-between items-center mb-6 mt-6">
        <div className="form-group form-check">
          <input
            type="checkbox"
            className="form-check-input appearance-none h-4 w-4 border border-gray-300 rounded-sm bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer"
            id="exampleCheck2"
          />
          <label
            className="form-check-label inline-block text-gray-800"
            htmlFor="exampleCheck2"
          >
            Remember me
          </label>
        </div>
        <a
          href="#"
          onClick={() => nagigateToOtp()}
          className="text-gray-800"
        >
          Forgot password?
        </a>
        <a
          href="#"
          onClick={() => nagigateToOtp()}
          className="text-gray-800"
        >
          Forgot password?
        </a>
      </div>
    </div>
  )
}

export default ForgotPass
